package types

type Azure struct {
  ResourceGroup  string `mapstructure:"resource_group"`
  ClusterName    string `mapstructure:"cluster_name"`
  SubscriptionId string `mapstructure:"subscription_id"`
}
