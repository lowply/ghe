package main

type Config struct {
	Enterprise struct {
		PrivateMode        bool   `json:"private_mode"`
		PublicPages        bool   `json:"public_pages"`
		SubdomainIsolation bool   `json:"subdomain_isolation"`
		SignupEnabled      bool   `json:"signup_enabled"`
		GithubHostname     string `json:"github_hostname"`
		// IdenticonsHost        interface{} `json:"identicons_host"`
		HTTPProxy           interface{} `json:"http_proxy"`
		HTTPNoproxy         interface{} `json:"http_noproxy"`
		AuthMode            string      `json:"auth_mode"`
		BuiltinAuthFallback bool        `json:"builtin_auth_fallback"`
		// ExpireSessions      interface{} `json:"expire_sessions"`
		AdminPassword string `json:"admin_password"`
		// ConfigurationID       interface{} `json:"configuration_id"`
		// ConfigurationRunCount interface{} `json:"configuration_run_count"`
		// Customer              struct {
		// 	Name          interface{} `json:"name"`
		// 	Email         interface{} `json:"email"`
		// 	UUID          interface{} `json:"uuid"`
		// 	SecretKeyData interface{} `json:"secret_key_data"`
		// 	PublicKeyData interface{} `json:"public_key_data"`
		// } `json:"customer"`
		// License struct {
		// 	Seats            interface{} `json:"seats"`
		// 	Evaluation       interface{} `json:"evaluation"`
		// 	Perpetual        interface{} `json:"perpetual"`
		// 	UnlimitedSeating interface{} `json:"unlimited_seating"`
		// 	SupportKey       interface{} `json:"support_key"`
		// 	SSHAllowed       interface{} `json:"ssh_allowed"`
		// 	ClusterSupport   interface{} `json:"cluster_support"`
		// 	ExpireAt         interface{} `json:"expire_at"`
		// } `json:"license"`
		GithubSsl struct {
			Enabled bool   `json:"enabled"`
			Cert    string `json:"cert"`
			Key     string `json:"key"`
			// TLSMode []string `json:"tls_mode"`
			Acme struct {
				Enabled bool `json:"enabled"`
				// AcceptTos      bool   `json:"accept_tos"`
				// ContactEmail   string `json:"contact_email"`
				// ValidationType string `json:"validation_type"`
				// Provider       string `json:"provider"`
			} `json:"acme"`
		} `json:"github_ssl"`
		Ldap struct {
			Host                      string        `json:"host"`
			Port                      int           `json:"port"`
			Base                      []string      `json:"base"`
			UID                       string        `json:"uid"`
			BindDn                    string        `json:"bind_dn"`
			Password                  string        `json:"password"`
			Method                    string        `json:"method"`
			SearchStrategy            string        `json:"search_strategy"`
			UserGroups                []interface{} `json:"user_groups"`
			AdminGroup                string        `json:"admin_group"`
			VirtualAttributeEnabled   bool          `json:"virtual_attribute_enabled"`
			RecursiveGroupSearch      bool          `json:"recursive_group_search"`
			PosixSupport              bool          `json:"posix_support"`
			UserSyncEmails            bool          `json:"user_sync_emails"`
			UserSyncKeys              bool          `json:"user_sync_keys"`
			UserSyncGpgKeys           bool          `json:"user_sync_gpg_keys"`
			UserSyncInterval          int           `json:"user_sync_interval"`
			TeamSyncInterval          int           `json:"team_sync_interval"`
			SyncEnabled               bool          `json:"sync_enabled"`
			ExternalAuthTokenRequired bool          `json:"external_auth_token_required"`
			VerifyCertificate         bool          `json:"verify_certificate"`
			// Reconciliation            struct {
			// 	User interface{} `json:"user"`
			// 	Org  interface{} `json:"org"`
			// } `json:"reconciliation"`
			Profile struct {
				UID    string      `json:"uid"`
				Name   interface{} `json:"name"`
				Mail   string      `json:"mail"`
				Key    string      `json:"key"`
				GpgKey interface{} `json:"gpg_key"`
			} `json:"profile"`
		} `json:"ldap"`
		// Cas struct {
		// 	URL interface{} `json:"url"`
		// } `json:"cas"`
		Saml struct {
			SsoURL      string `json:"sso_url"`
			Certificate string `json:"certificate"`
			// CertificatePath    string      `json:"certificate_path"`
			Issuer             string `json:"issuer"`
			NameIDFormat       string `json:"name_id_format"`
			IdpInitiatedSso    bool   `json:"idp_initiated_sso"`
			DisableAdminDemote bool   `json:"disable_admin_demote"`
			// SignatureMethod    string      `json:"signature_method"`
			// DigestMethod       string      `json:"digest_method"`
			// UsernameAttribute interface{} `json:"username_attribute"`
			// FullNameAttribute string      `json:"full_name_attribute"`
			// EmailsAttribute   string      `json:"emails_attribute"`
			// SSHKeysAttribute  string      `json:"ssh_keys_attribute"`
			// GpgKeysAttribute  string      `json:"gpg_keys_attribute"`
		} `json:"saml"`
		GithubOauth interface{} `json:"github_oauth"`
		SMTP        struct {
			Enabled bool   `json:"enabled"`
			Address string `json:"address"`
			// Authentication interface{} `json:"authentication"`
			Port int `json:"port"`
			// Domain                  interface{} `json:"domain"`
			// Username                interface{} `json:"username"`
			// UserName                interface{} `json:"user_name"`
			// Password                interface{} `json:"password"`
			SupportAddress     string `json:"support_address"`
			SupportAddressType string `json:"support_address_type"`
			NoreplyAddress     string `json:"noreply_address"`
			// DiscardToNoreplyAddress bool   `json:"discard_to_noreply_address"`
		} `json:"smtp"`
		// Ntp struct {
		// 	PrimaryServer   string `json:"primary_server"`
		// 	SecondaryServer string `json:"secondary_server"`
		// } `json:"ntp"`
		// Timezone interface{} `json:"timezone"`
		Snmp struct {
			Enabled   bool          `json:"enabled"`
			Version   int           `json:"version"`
			Community string        `json:"community"`
			Users     []interface{} `json:"users"`
		} `json:"snmp"`
		Syslog struct {
			Enabled      bool   `json:"enabled"`
			Server       string `json:"server"`
			ProtocolName string `json:"protocol_name"`
			TLSEnabled   bool   `json:"tls_enabled"`
			Cert         string `json:"cert"`
		} `json:"syslog"`
		// Assets interface{} `json:"assets"`
		Pages struct {
			Enabled bool `json:"enabled"`
		} `json:"pages"`
		Collectd struct {
			Enabled    bool        `json:"enabled"`
			Server     interface{} `json:"server"`
			Port       int         `json:"port"`
			Encryption interface{} `json:"encryption"`
			Username   interface{} `json:"username"`
			Password   interface{} `json:"password"`
		} `json:"collectd"`
		// Mapping struct {
		// 	Enabled    bool        `json:"enabled"`
		// 	Tileserver interface{} `json:"tileserver"`
		// 	Basemap    string      `json:"basemap"`
		// 	Token      interface{} `json:"token"`
		// } `json:"mapping"`
		LoadBalancer struct {
			HTTPForward   bool `json:"http_forward"`
			ProxyProtocol bool `json:"proxy_protocol"`
		} `json:"load_balancer"`
		// AbuseRateLimiting struct {
		// 	Enabled                  bool `json:"enabled"`
		// 	RequestsPerMinute        int  `json:"requests_per_minute"`
		// 	CPUMillisPerMinute       int  `json:"cpu_millis_per_minute"`
		// 	SearchCPUMillisPerMinute int  `json:"search_cpu_millis_per_minute"`
		// } `json:"abuse_rate_limiting"`
		// APIRateLimiting struct {
		// 	Enabled                         bool `json:"enabled"`
		// 	UnauthenticatedRateLimit        int  `json:"unauthenticated_rate_limit"`
		// 	DefaultRateLimit                int  `json:"default_rate_limit"`
		// 	SearchUnauthenticatedRateLimit  int  `json:"search_unauthenticated_rate_limit"`
		// 	SearchDefaultRateLimit          int  `json:"search_default_rate_limit"`
		// 	LfsUnauthenticatedRateLimit     int  `json:"lfs_unauthenticated_rate_limit"`
		// 	LfsDefaultRateLimit             int  `json:"lfs_default_rate_limit"`
		// 	GraphqlUnauthenticatedRateLimit int  `json:"graphql_unauthenticated_rate_limit"`
		// 	GraphqlDefaultRateLimit         int  `json:"graphql_default_rate_limit"`
		// } `json:"api_rate_limiting"`
		// Governor struct {
		// 	QuotasEnabled bool        `json:"quotas_enabled"`
		// 	LimitUser     interface{} `json:"limit_user"`
		// 	LimitNetwork  interface{} `json:"limit_network"`
		// } `json:"governor"`
		// Applications struct {
		// 	AvatarsMaxAge int `json:"avatars_max_age"`
		// } `json:"applications"`
	} `json:"enterprise"`
	// RunList []string `json:"run_list"`
}
