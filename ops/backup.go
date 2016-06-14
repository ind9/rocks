package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

// Current is the file which terminates recursive loops
const Current = "CURRENT"

var backup = &cobra.Command{
	Use:   "backup",
	Short: "Backs up rocksdb stores",
	Long:  "Backs up rocksdb stores",
	Run:   AttachHandler(backupDatabase),
}

func backupDatabase(args []string) error {
	if source == "" {
		return fmt.Errorf("--src was not set")
	}
	if destination == "" {
		return fmt.Errorf("--dest was not set")
	}
	if recursive {
		walkSourceDir(source, destination)
		return nil
	}
	log.Printf("Trying to create backup from %s to %s\n", source, destination)
	return DoBackup(source, destination)
}

func walkSourceDir(source, destination string) {
	filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {

		if info.Name() == Current {
			dbLoc := filepath.Dir(path)
			dbBackupLoc := filepath.Join(destination, dbLoc)
			log.Printf("Backup at %s, would be stored to %s\n", dbLoc, dbBackupLoc)

			dbRelative, err := filepath.Rel(source, dbLoc)
			if err != nil {
				log.Print(err)
				return err
			}
			log.Printf("Backup is created for %s rocks store", dbRelative)
			if err = os.MkdirAll(dbBackupLoc, os.ModePerm); err != nil {
				log.Print(err)
				return err
			}

			if err = DoBackup(dbLoc, dbBackupLoc); err != nil {
				log.Print(err)
				return err
			}
			return filepath.SkipDir
		}
		return walkErr
	})
}

// DoBackup triggers a backup from the source
func DoBackup(source, destination string) error {

	opts := gorocksdb.NewDefaultOptions()
	db, err := gorocksdb.OpenDb(opts, source)
	if err != nil {
		return err
	}

	backup, err := gorocksdb.OpenBackupEngine(opts, destination)
	if err != nil {
		return err
	}
	return backup.CreateNewBackup(db)
}

func init() {
	Rocks.AddCommand(backup)

	backup.PersistentFlags().StringVar(&source, "src", "", "Backup from")
	backup.PersistentFlags().StringVar(&destination, "dest", "", "Backup to")
	backup.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to backup in recursive fashion from src to dest")
}
