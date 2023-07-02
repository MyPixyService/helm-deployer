#!/bin/bash

# Check if the script is being run as root
if [ "$EUID" -ne 0 ]; then
    echo "This script must be run with sudo." >&2
    echo "You will be prompted for your password." >&2
    sudo "$0" "$@"
    exit $?
fi

echo "Running with root privileges..."

echo "Downloading latest binary..."
wget -q --show-progress -O helm-deployer https://github.com/MyPixyService/helm-deployer/releases/download/v0.1.0-alpha/linux-helm-deployer-v0.1.0-alpha

echo "Moving binary to installation directory..."
mv helm-deployer /usr/local/bin/

echo "Adding execution bit..."
chmod +x /usr/local/bin/helm-deployer

echo -e "\nHelm Deployer has successfully been installed on your system."
echo "Usage: helm-deployer -f <inputfile.yaml> -k <kubeconfig> [-c <releaseName>] [-uninstall]"
echo "For more information visit: https://github.com/MyPixyService/helm-deployer"
