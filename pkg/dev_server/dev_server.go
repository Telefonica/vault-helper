package dev_server

import (
	"github.com/Sirupsen/logrus"

	"github.com/jetstack-experimental/vault-helper/pkg/kubernetes"
	"github.com/jetstack-experimental/vault-helper/pkg/testing/vault_dev"
)

type DevVault struct {
	Vault      *vault_dev.VaultDev
	Kubernetes *kubernetes.Kubernetes
	Log        *logrus.Entry
}

func New(logger *logrus.Entry) *DevVault {

	vault := vault_dev.New()
	if err := vault.Start(); err != nil {
		logrus.Fatalf("unable to initialise dev vault: %s", err)
	}

	v := &DevVault{
		Vault: vault,
		Log:   logger,
	}

	return v
}