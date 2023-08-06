package civo

import (
	"strings"

	"github.com/kubesimplify/ksctl/api/resources"
	"github.com/kubesimplify/ksctl/api/utils"
)

// DelSSHKeyPair implements resources.CloudInfrastructure.
func (obj *CivoProvider) DelSSHKeyPair(storage resources.StorageInfrastructure) error {
	if len(civoCloudState.SSHID) == 0 {
		storage.Logger().Success("[skip] ssh keypair already deleted")
		return nil
	}

	_, err := civoClient.DeleteSSHKey(civoCloudState.SSHID)
	if err != nil {
		return err
	}
	path := generatePath(utils.CLUSTER_PATH, clusterType, clusterDirName, STATE_FILE_NAME)

	storage.Logger().Success("[civo] ssh keypair deleted", civoCloudState.SSHID)

	civoCloudState.SSHID = ""
	civoCloudState.SSHPrivateKeyLoc = ""
	civoCloudState.SSHUser = ""

	return saveStateHelper(storage, path)
}

// CreateUploadSSHKeyPair implements resources.CloudInfrastructure.
func (obj *CivoProvider) CreateUploadSSHKeyPair(storage resources.StorageInfrastructure) error {
	if len(civoCloudState.SSHID) != 0 {
		storage.Logger().Success("[skip] ssh keypair already uploaded")
		return nil
	}

	keyPairToUpload, err := utils.CreateSSHKeyPair(storage, "civo", clusterDirName)
	if err != nil {
		return err
	}
	if err := obj.uploadSSH(storage, obj.Metadata.ResName, keyPairToUpload); err != nil {
		return err
	}
	storage.Logger().Success("[civo] ssh keypair created and uploaded", obj.Metadata.ResName)
	return nil
}

func (obj *CivoProvider) uploadSSH(storage resources.StorageInfrastructure, resName, pubKey string) error {
	sshResp, err := civoClient.NewSSHKey(strings.ToLower(resName), pubKey)
	if err != nil {
		return err
	}

	// NOTE: state for the ssh
	civoCloudState.SSHID = sshResp.ID
	civoCloudState.SSHUser = "root"
	civoCloudState.SSHPrivateKeyLoc = utils.GetPath(utils.SSH_PATH, "civo", clusterType, clusterDirName)

	path := generatePath(utils.CLUSTER_PATH, clusterType, clusterDirName, STATE_FILE_NAME)

	return saveStateHelper(storage, path)
}
