package azure

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Humalect/humalect-core/agent/constants"
	"github.com/Humalect/humalect-core/agent/services/k8s"
	"k8s.io/client-go/kubernetes"
)

func CreateAcrSecret(params constants.ParamsConfig, clientSet *kubernetes.Clientset) (string, error) {
	var acrCredentials constants.AcrCredentials
	_ = json.Unmarshal([]byte(params.AcrCredentials), &acrCredentials)
	azureCreds, err := FetchAcrCreds(acrCredentials.ManagementScopeToken, acrCredentials.RegistryName,
		acrCredentials.SubscriptionId, acrCredentials.ResourceGroupName)
	if err != nil {
		log.Fatalf("Error getting Azure ACR creds: %v", err)
		return "", err
	}

	secretData := fmt.Sprintf(`{  
			"auths": {  
				"%s.azurecr.io": {  
					"username": "%s",  
					"password": "%s"  
				}  
			}  
		}`, acrCredentials.RegistryName, azureCreds.Username, azureCreds.Password)
	return k8s.CreateSecret(map[string]string{
		".dockerconfigjson": secretData,
	}, params, clientSet, "humalect")
}
