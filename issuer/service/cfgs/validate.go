package cfgs

import "fmt"

func validateConfig(cfg *IssuerConfig) error {
	if len(cfg.LogLevel) < 4 {
		return fmt.Errorf(`the config parameter "\log_level"\ wasn't specified'`)
	}

	if len(cfg.DBFilePath) == 0 {
		return fmt.Errorf(`the config parameter "\db_file_path"\ wasn't specified'`)
	}

	if len(cfg.LocalUrl) < 5 { // at least specifying port ":xxxx"
		return fmt.Errorf(`the config parameter "\local_url"\ wasn't specified'`)
	}

	if len(cfg.PublicUrl) < 4 { // at least specifying port ":xxxx"
		return fmt.Errorf(`the config parameter "\public_url"\ wasn't specified'`)
	}

	if len(cfg.NodeRpcUrl) == 0 {
		return fmt.Errorf(`the config parameter "\node_rpc_url"\ wasn't specified'`)
	}

	if len(cfg.PublishingContractAddress) < 32 {
		return fmt.Errorf(`the config parameter "\publishing_contract_address"\ wasn't specified'`)
	}

	if len(cfg.PublishingPrivateKey) < 32 {
		return fmt.Errorf(`the config parameter "\publishing_private_key"\ wasn't specified'`)
	}

	if len(cfg.CircuitsDir) == 0 {
		return fmt.Errorf(`the config parameter "\circuits_dir"\ wasn't specified'`)
	}

	if len(cfg.IpfsUrl) == 0 {
		return fmt.Errorf(`the config parameter "\ipfs_url"\ wasn't specified'`)
	}

	return nil
}
