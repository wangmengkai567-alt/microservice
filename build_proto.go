package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	pwd, _ := os.Getwd()
	protocPath := filepath.Join(pwd, "protoc_bin", "bin", "protoc")

	// Run protoc
	os.Setenv("PATH", os.Getenv("PATH")+";"+filepath.Join(pwd, "protoc_bin", "bin")+";E:\\go_env\\bin")

	cmd1 := exec.Command(protocPath, "--go_out=.", "--go_opt=paths=source_relative", "--go-grpc_out=.", "--go-grpc_opt=paths=source_relative", "proto/product.proto")
	cmd1.Dir = "product-service"
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	if err := cmd1.Run(); err != nil {
		log.Fatal(err)
	}

	cmd2 := exec.Command(protocPath, "--go_out=.", "--go_opt=paths=source_relative", "--go-grpc_out=.", "--go-grpc_opt=paths=source_relative", "proto/product/product.proto")
	cmd2.Dir = "api_gateway"
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	if err := cmd2.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully compiled proto files.")
}

func unzip(src, dest string) {
	r, err := zip.OpenReader(src)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			log.Fatal(err)
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatal(err)
		}
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
