package subscription

import (
	"chs_cloud_general/internal/global_var"
	"time"
)

type SubscriptionStruct struct {
	SubscriptionId uint64
	Subdomain      string
	Domain         string
	StartDate      time.Time
	EndDate        time.Time
	AddOn          []int
	IsActive       bool
}

type DataSubscriptionStruct struct {
	SubscriptionId uint64
	Subdomain      string
	Domain         string
	StartDate      time.Time
	EndDate        time.Time
	IsActive       bool
}

type AvailableModuleStruct struct {
	ID   int
	Code string
	Name string
}

var AvailableModule = []AvailableModuleStruct{
	0: {
		ID:   0,
		Name: "Front Desk",
		Code: global_var.SystemCode.Hotel,
	},
	1: {
		ID:   1,
		Name: "Point of Sales",
		Code: global_var.SystemCode.Pos,
	},
	2: {
		ID:   2,
		Name: "Banquet",
		Code: global_var.SystemCode.Banquet,
	},
	3: {
		ID:   3,
		Name: "Accounting",
		Code: global_var.SystemCode.Accounting,
	},
	4: {
		ID:   4,
		Name: "Inventory & Assets",
		Code: global_var.SystemCode.Asset,
	},
	5: {
		ID:   5,
		Name: "Report",
		Code: global_var.SystemCode.Report,
	},
	6: {
		ID:   6,
		Name: "Tools",
		Code: global_var.SystemCode.Tools,
	},
}
var ActiveSubscriptions = make(map[string]SubscriptionStruct)

func LoadDataSubscription() {
	var Subscription []DataSubscriptionStruct

	// Generate subscriptions by subdomain/domain
	global_var.DBMain.Table("company").Select("subscription_id,subdomain,domain,subscription.start_date,subscription.end_date").
		Joins("LEFT JOIN company_database ON company.company_id = company_database.company_id").
		Joins("LEFT JOIN subscription ON company.company_id = subscription.company_id").
		Scan(&Subscription)

	for _, subscription := range Subscription {
		var AddonID []int
		global_var.DBMain.Table("subscription_addon").Select("IF(addon_id = 99,0,IFNULL(addon_id, 0)) AS AddonID").Where("subscription_id = ?", subscription.SubscriptionId).
			Scan(&AddonID)

		// Only front desk
		if len(AddonID) <= 0 {
			AddonID = []int{0}
		}

		if subscription.Subdomain != "" {
			key := subscription.Subdomain

			ActiveSubscriptions[key] = SubscriptionStruct{
				SubscriptionId: subscription.SubscriptionId,
				Subdomain:      subscription.Subdomain,
				Domain:         subscription.Domain,
				StartDate:      subscription.StartDate,
				EndDate:        subscription.EndDate,
				AddOn:          AddonID,
				IsActive:       subscription.StartDate.Before(time.Now()) && subscription.EndDate.After(time.Now()),
			}
		}
		if subscription.Domain != "" {
			key := subscription.Domain
			ActiveSubscriptions[key] = SubscriptionStruct{
				SubscriptionId: subscription.SubscriptionId,
				Subdomain:      subscription.Subdomain,
				Domain:         subscription.Domain,
				StartDate:      subscription.StartDate,
				EndDate:        subscription.EndDate,
				AddOn:          AddonID,
				IsActive:       subscription.StartDate.Before(time.Now()) && subscription.EndDate.After(time.Now()),
			}
		}
	}
}
