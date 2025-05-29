#!/bin/bash

echo "=== CHECKING FOR POTENTIALLY MISSING COMMON AZURE RESOURCES ==="
echo ""

# List of common/newer Azure resources that might be missing
common_resources=(
    "azurerm_monitor_workspace"
    "azurerm_grafana"
    "azurerm_chaos_studio_experiment" 
    "azurerm_chaos_studio_target"
    "azurerm_container_app_job"
    "azurerm_elastic_san"
    "azurerm_elastic_san_volume_group"
    "azurerm_elastic_san_volume"
    "azurerm_nginx_certificate"
    "azurerm_nginx_configuration"
    "azurerm_orbital_contact_profile"
    "azurerm_orbital_spacecraft"
    "azurerm_palo_alto_local_rulestack"
    "azurerm_palo_alto_virtual_network_appliance"
    "azurerm_voice_services_communications_gateway"
    "azurerm_app_service_certificate_binding"
    "azurerm_app_service_certificate_order"
    "azurerm_app_service_hybrid_connection"
    "azurerm_app_service_managed_certificate"
    "azurerm_app_service_public_certificate"
    "azurerm_app_service_slot_custom_hostname_binding"
    "azurerm_app_service_slot_virtual_network_swift_connection"
    "azurerm_app_service_source_control_slot"
    "azurerm_app_service_source_control_token"
    "azurerm_app_service_virtual_network_swift_connection"
    "azurerm_backup_container_storage_account"
    "azurerm_backup_policy_file_share"
    "azurerm_backup_policy_vm_workload"
    "azurerm_backup_protected_file_share"
    "azurerm_backup_protected_vm"
    "azurerm_sentinel_alert_rule"
    "azurerm_sentinel_automation_rule"
    "azurerm_sentinel_data_connector"
    "azurerm_sentinel_watchlist"
    "azurerm_virtual_machine_extension"
    "azurerm_virtual_machine_scale_set_extension"
    "azurerm_windows_function_app"
    "azurerm_linux_function_app"
    "azurerm_function_app_active_slot"
    "azurerm_function_app_function"
    "azurerm_function_app_hybrid_connection"
)

echo "Checking for common resources that might be missing..."
echo ""

missing_common=0
for resource in "${common_resources[@]}"; do
    if ! jq -e ".\"$resource\"" /workspaces/terraform-provider-restful/internal/provider/resource_definitions.json > /dev/null 2>&1; then
        echo "POTENTIALLY MISSING: $resource"
        ((missing_common++))
    fi
done

if [ $missing_common -eq 0 ]; then
    echo "All checked common resources are already present in the JSON."
fi

echo ""
echo "=== PROPOSED SLUGS FOR ANY MISSING RESOURCES ==="
echo ""

# Function to suggest a slug based on resource name
suggest_slug() {
    local resource_name="$1"
    
    # Remove azurerm_ prefix
    local base_name="${resource_name#azurerm_}"
    
    # Create abbreviated slug based on common patterns
    case "$base_name" in
        *monitor_workspace*) echo "amws" ;;
        *grafana*) echo "grf" ;;
        *chaos_studio_experiment*) echo "chse" ;;
        *chaos_studio_target*) echo "chst" ;;
        *container_app_job*) echo "caj" ;;
        *elastic_san*) echo "esan" ;;
        *elastic_san_volume_group*) echo "esanvg" ;;
        *elastic_san_volume*) echo "esanv" ;;
        *nginx_certificate*) echo "ngxcrt" ;;
        *nginx_configuration*) echo "ngxcfg" ;;
        *orbital_contact_profile*) echo "orcp" ;;
        *orbital_spacecraft*) echo "orsc" ;;
        *palo_alto_local_rulestack*) echo "palrs" ;;
        *palo_alto_virtual_network_appliance*) echo "palvna" ;;
        *voice_services_communications_gateway*) echo "vscgw" ;;
        *app_service_certificate_binding*) echo "appcb" ;;
        *app_service_certificate_order*) echo "appco" ;;
        *app_service_hybrid_connection*) echo "apphc" ;;
        *app_service_managed_certificate*) echo "appmc" ;;
        *app_service_public_certificate*) echo "apppc" ;;
        *backup_container_storage_account*) echo "bcsa" ;;
        *backup_policy_file_share*) echo "bpfs" ;;
        *backup_policy_vm_workload*) echo "bpvw" ;;
        *backup_protected_file_share*) echo "bpfs2" ;;
        *backup_protected_vm*) echo "bpvm" ;;
        *sentinel_alert_rule*) echo "senar" ;;
        *sentinel_automation_rule*) echo "senaut" ;;
        *sentinel_data_connector*) echo "sendc" ;;
        *sentinel_watchlist*) echo "senwl" ;;
        *virtual_machine_extension*) echo "vmext" ;;
        *virtual_machine_scale_set_extension*) echo "vmssext" ;;
        *windows_function_app*) echo "wfa" ;;
        *linux_function_app*) echo "lfa" ;;
        *function_app_active_slot*) echo "faas" ;;
        *function_app_function*) echo "faf" ;;
        *function_app_hybrid_connection*) echo "fahc" ;;
        *) 
            # Generic slug generation - take first letters of major words
            echo "$base_name" | sed 's/_/ /g' | awk '{for(i=1;i<=NF;i++) printf substr($i,1,1); print ""}' | tr '[:upper:]' '[:lower:]' | head -c 6
            ;;
    esac
}

# Show suggested slugs for missing resources
for resource in "${common_resources[@]}"; do
    if ! jq -e ".\"$resource\"" /workspaces/terraform-provider-restful/internal/provider/resource_definitions.json > /dev/null 2>&1; then
        suggested_slug=$(suggest_slug "$resource")
        echo "$resource -> suggested slug: $suggested_slug"
    fi
done

echo ""
echo "=== SUMMARY ==="
echo "Total potentially missing common resources: $missing_common"
