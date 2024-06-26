#!/bin/bash

# Golangのバージョンを指定
GOVERSION="1.20.5"

# OSチェック
OS="$(uname)"
case $OS in
    "Darwin")
        echo "MacOS detected."
        if ! command -v brew &> /dev/null; then
            echo "Homebrew is not installed. Installing Homebrew..."
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        fi
        echo "Installing Golang..."
        brew install go@$GOVERSION
        ;;
    "Linux")
        echo "Linux detected."
        if ! command -v apt-get &> /dev/null; then
            echo "This script supports only Debian-based distributions."
            exit 1
        fi
        echo "Installing Golang..."
        sudo apt-get update
        sudo apt-get install -y wget tar
        wget https://dl.google.com/go/go$GOVERSION.linux-amd64.tar.gz
        sudo tar -C /usr/local -xzf go$GOVERSION.linux-amd64.tar.gz
        export PATH=$PATH:/usr/local/go/bin
        ;;
    "MINGW"*|"MSYS_NT"*)
        echo "Windows detected."
        echo "Please install Golang manually from https://golang.org/dl/ and add it to your PATH."
        read -p "Press Enter to continue after installing Golang..."
        ;;
    *)
        echo "Unsupported OS detected: $OS"
        exit 1
        ;;
esac

# Golangのバージョン確認
go version

# 必要なGolangモジュールのインストール
echo "Installing required Golang modules..."
go mod tidy

echo "Setup complete."

