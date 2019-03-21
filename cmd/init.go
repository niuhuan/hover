package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a flutter project to use go-flutter",
	Run: func(cmd *cobra.Command, args []string) {
		assertInFlutterProject()

		err := os.Mkdir("desktop", 0775)
		if err != nil {
			if os.IsExist(err) {
				fmt.Println("A file or directory named `desktop` already exists. Cannot continue init.")
				os.Exit(1)
			}
		}

		desktopCmdPath := filepath.Join("desktop", "cmd")
		err = os.Mkdir(desktopCmdPath, 0775)
		if err != nil {
			fmt.Printf("Failed to create `%s`: %v\n", desktopCmdPath, err)
			os.Exit(1)
		}

		desktopAssetsPath := filepath.Join("desktop", "assets")
		err = os.Mkdir(desktopAssetsPath, 0775)
		if err != nil {
			fmt.Printf("Failed to create `%s`: %v\n", desktopAssetsPath, err)
			os.Exit(1)
		}

		copyAsset("app/main.go", filepath.Join(desktopCmdPath, "main.go"))
		copyAsset("app/options.go", filepath.Join(desktopCmdPath, "options.go"))
		copyAsset("app/logo.png", filepath.Join(desktopAssetsPath, "logo.png"))
		copyAsset("app/gitignore", filepath.Join("desktop", ".gitignore"))

		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Failed to get working dir: %v\n", err)
			os.Exit(1)
		}

		cmdGoModInit := exec.Command(goBin, "mod", "init")
		cmdGoModInit.Dir = filepath.Join(wd, "desktop")
		cmdGoModInit.Env = append(os.Environ(),
			"GO111MODULE=on",
		)
		cmdGoModInit.Stderr = os.Stderr
		cmdGoModInit.Stdout = os.Stdout
		err = cmdGoModInit.Run()
		if err != nil {
			fmt.Printf("Go mod init failed: %v\n", err)
			os.Exit(1)
		}

		cmdGoModTidy := exec.Command(goBin, "mod", "tidy")
		cmdGoModTidy.Dir = filepath.Join(wd, "desktop")
		fmt.Println(cmdGoModTidy.Dir)
		cmdGoModTidy.Env = append(os.Environ(),
			"GO111MODULE=on",
		)
		cmdGoModTidy.Stderr = os.Stderr
		cmdGoModTidy.Stdout = os.Stdout
		err = cmdGoModTidy.Run()
		if err != nil {
			fmt.Printf("Go mod tidy failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func copyAsset(boxed, to string) {
	file, err := os.Create(to)
	if err != nil {
		fmt.Printf("Failed to create %s: %v\n", to, err)
		os.Exit(1)
	}
	defer file.Close()
	boxedFile, err := assetsBox.Open(boxed)
	if err != nil {
		fmt.Printf("Failed to find boxed file %s: %v\n", boxed, err)
		os.Exit(1)
	}
	defer boxedFile.Close()
	_, err = io.Copy(file, boxedFile)
	if err != nil {
		fmt.Printf("Failed to write file %s: %v\n", to, err)
		os.Exit(1)
	}
}
