package backup

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ashwanthkumar/golang-utils/worker"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

var backupSource string
var backupDestination string
var backupThreads int
var recursive bool

var backup = &cobra.Command{
	Use:   "backup",
	Short: "Backs up rocksdb stores",
	Long:  "Backs up rocksdb stores",
	Run:   AttachHandler(backupDatabase),
}

// Work struct contains source and destination for backup
type Work struct {
	Source      string
	Destination string
}

func backupDatabase(args []string) error {
	if backupSource == "" {
		return fmt.Errorf("--src was not set")
	}
	if backupDestination == "" {
		return fmt.Errorf("--dest was not set")
	}
	if recursive {
		return DoRecursiveBackup(backupSource, backupDestination, backupThreads)
	}
	return DoBackup(backupSource, backupDestination)
}

// DoRecursiveBackup recursively takes a rocksdb backup keeping the folder structure intact as in source
func DoRecursiveBackup(source, destination string, threads int) error {

	workerPool := worker.Pool{
		MaxWorkers: threads,
		Op: func(request worker.Request) error {
			work := request.(Work)
			return DoBackup(work.Source, work.Destination)
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == Current {
			dbLoc := filepath.Dir(path)

			dbRelative, err := filepath.Rel(source, dbLoc)
			if err != nil {
				log.Print(err)
				return err
			}

			dbBackupLoc := filepath.Join(destination, dbRelative)
			if err = os.MkdirAll(dbBackupLoc, os.ModePerm); err != nil {
				log.Print(err)
				return err
			}

			work := BackupWork{
				Source:      dbLoc,
				Destination: dbBackupLoc,
			}
			workerPool.AddWork(work)
			return filepath.SkipDir
		}
		return walkErr
	})

	var result error
	if errFromWorkers := workerPool.Join(); errFromWorkers != nil {
		result = multierror.Append(result, errFromWorkers)
	}

	if err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// DoBackup triggers a backup from the source
func DoBackup(source, destination string) error {
	log.Printf("Trying to create backup from %s to %s\n", source, destination)

	opts := gorocksdb.NewDefaultOptions()
	db, err := gorocksdb.OpenDb(opts, source)
	if err != nil {
		return err
	}

	backup, err := gorocksdb.OpenBackupEngine(opts, destination)
	if err != nil {
		return err
	}
	err = backup.CreateNewBackup(db)
	db.Close()
	backup.Close()
	log.Printf("Backup from %s to %s completed\n", source, destination)
	return err
}

func init() {
	Rocks.AddCommand(backup)

	backup.PersistentFlags().StringVar(&backupSource, "src", "", "Backup from")
	backup.PersistentFlags().StringVar(&backupDestination, "dest", "", "Backup to")
	backup.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to backup in recursive fashion from src to dest")
	backup.PersistentFlags().IntVar(&backupThreads, "threads", 2*runtime.NumCPU(), "Number of threads to do backup")
}
