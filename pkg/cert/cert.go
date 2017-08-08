package cert

import (
	"encoding/pem"
	"fmt"
	"os"
	"os/user"
	"strconv"

	"github.com/Sirupsen/logrus"
	vault "github.com/hashicorp/vault/api"
	"github.com/jetstack-experimental/vault-helper/pkg/instanceToken"
)

type Cert struct {
	role        string
	commonName  string
	destination string
	bitSize     int
	pemSize     int
	keyType     string
	ipSans      []string
	sanHosts    []string
	owner       string
	group       string
	data        *pem.Block
	configPath  string

	vaultClient *vault.Client
	Log         *logrus.Entry
}

func (c *Cert) RunCert() error {
	if err := c.EnsureKey(); err != nil {
		return fmt.Errorf("error ensuring key: %s", err)
	}

	if err := c.TokenRenew(); err != nil {
		return fmt.Errorf("error renewing tokens: %s", err)
	}

	if err := c.RequestCertificate(); err != nil {
		return fmt.Errorf("error requesting certificate: %s", err)
	}
	return nil
}

func (c *Cert) TokenRenew() error {
	i := instanceToken.New(c.vaultClient, c.Log)
	i.SetRole(c.Role())
	i.SetVaultConfigPath(c.VaultConfigPath())

	return i.TokenRenewRun()
}

func (c *Cert) DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error removing file at '%s':  %s", path, err)
	}

	return nil
}

func (c *Cert) WritePermissions(path string, perm os.FileMode) error {

	if err := os.Chmod(path, perm); err != nil {
		return fmt.Errorf("failed to change permissons of file '%s' to 0600: %s", path, err)
	}

	usr, err := user.Lookup(c.Owner())
	if err != nil {
		return fmt.Errorf("failed to find user '%s' on system: %s", c.Owner(), err)
	}
	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return fmt.Errorf("failed to convert user uid '%s' (string) to (int): %s", usr.Uid, err)
	}
	grp, err := user.LookupGroup(c.Group())
	if err != nil {
		return fmt.Errorf("failed to find group '%s' on system: %s", c.Group(), err)
	}
	gid, err := strconv.Atoi(grp.Gid)
	if err != nil {
		return fmt.Errorf("failed to convert group gid '%s' (string) to (int): %s", grp.Gid, err)
	}

	if err := os.Chown(path, uid, gid); err != nil {
		return fmt.Errorf("failed to change group and owner of file '%s' to usr:'%s' grp:'%s': %s", path, c.owner, c.group, err)
	}

	c.Log.Debugf("Set permissons on file: %s", path)

	return nil
}

func New(vaultClient *vault.Client, logger *logrus.Entry) *Cert {
	c := &Cert{
		role:        "",
		commonName:  "",
		destination: "",
		bitSize:     2048,
		keyType:     "RSA",
		ipSans:      []string{},
		sanHosts:    []string{},
		owner:       "",
		group:       "",
	}

	if vaultClient != nil {
		c.vaultClient = vaultClient
	}
	if logger != nil {
		c.Log = logger
	}

	return c
}

func (c *Cert) SetRole(role string) {
	c.role = role
}
func (c *Cert) Role() string {
	return c.role
}

func (c *Cert) SetCommonName(name string) {
	c.commonName = name
}
func (c *Cert) CommonName() string {
	return c.commonName
}

func (c *Cert) SetDestination(destination string) {
	c.destination = destination
}
func (c *Cert) Destination() string {
	return c.destination
}

func (c *Cert) SetBitSize(size int) {
	c.bitSize = size
}
func (c *Cert) BitSize() int {
	return c.bitSize
}

func (c *Cert) SetPemSize(size int) {
	c.pemSize = size
}
func (c *Cert) PemSize() int {
	return c.pemSize
}

func (c *Cert) SetKeyType(keyType string) {
	c.keyType = keyType
}
func (c *Cert) KeyType() string {
	return c.keyType
}

func (c *Cert) SetIPSans(ips []string) {
	c.ipSans = ips
}
func (c *Cert) IPSans() []string {
	return c.ipSans
}

func (c *Cert) SetSanHosts(hosts []string) {
	c.sanHosts = hosts
}
func (c *Cert) SanHosts() []string {
	return c.sanHosts
}

func (c *Cert) SetOwner(owner string) {
	c.owner = owner
}
func (c *Cert) Owner() string {
	return c.owner
}

func (c *Cert) SetGroup(group string) {
	c.group = group
}
func (c *Cert) Group() string {
	return c.group
}

func (c *Cert) SetData(data *pem.Block) {
	c.data = data
}
func (c *Cert) Data() *pem.Block {
	return c.data
}

func (c *Cert) SetVaultConfigPath(path string) {
	c.configPath = path
}
func (c *Cert) VaultConfigPath() string {
	return c.configPath
}
