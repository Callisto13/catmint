package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Callisto13/pugo/pure1"
	"github.com/inhies/go-bytesize"
)

type config struct {
	subName     string
	licenseName string
	appID       string
	pemFile     string
}

func main() {
	cfg := &config{}

	flag.StringVar(&cfg.appID, "app-id", "", "The application ID you generated in Pure1")
	flag.StringVar(&cfg.pemFile, "private-key", "", "The path to your private rsa key")
	flag.StringVar(&cfg.subName, "sub", "", "The id (contract id) of the subscription to get data for")
	flag.StringVar(&cfg.licenseName, "license", "", "The name of the subcription license to get data for")

	flag.Parse()

	if err := validate(cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	privateKey, err := os.ReadFile(cfg.pemFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	restVersion := "1.2.b"

	client, err := pure1.NewClient(cfg.appID, privateKey, restVersion)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfg.licenseName != "" {
		if err := subscriptionLicenses(client, cfg); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	if err := subscriptionInfo(client, cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validate(cfg *config) error {
	if cfg.appID == "" {
		return errors.New("app-id required")
	}

	if cfg.pemFile == "" {
		return errors.New("private-key required")
	}

	return nil
}

func subscriptionInfo(client *pure1.Client, cfg *config) error {
	params := map[string]string{
		"subscription_names": fmt.Sprintf("'%s'", cfg.subName),
	}

	assets, err := client.Subscriptions.GetSubscriptionAssets(params)
	if err != nil {
		return err
	}

	for _, asset := range assets {
		cfg.licenseName = asset.License.Name // not a good idea but oh well!
		if err := subscriptionLicenses(client, cfg); err != nil {
			return err
		}
	}

	return nil
}

func subscriptionLicenses(client *pure1.Client, cfg *config) error {
	params := map[string]string{
		"names": fmt.Sprintf("'%s'", cfg.licenseName),
	}

	subs, err := client.Subscriptions.GetSubscriptionLicenses(params)
	if err != nil {
		return err
	}

	for _, sub := range subs {
		rtb, err := bytesize.Parse(fmt.Sprintf("%.0f %s", sub.Reservation.Data, sub.Reservation.Unit))
		if err != nil {
			return err
		}

		etb, err := bytesize.Parse(fmt.Sprintf("%.0f %s", sub.Usage.Data, sub.Usage.Unit))
		if err != nil {
			return err
		}

		fmt.Printf("Data for License '%s' in Subscription '%s'\n", cfg.licenseName, cfg.subName)
		fmt.Printf("Effective used: %s out of %s\n", etb, rtb)

		fmt.Println("")
	}

	return nil
}

// func subscriptionAssets(client *pure1.Client, cfg *config) {
// 	params := map[string]string{
// 		"subscription_names": fmt.Sprintf("'%s'", cfg.subName),
// 	}

// 	subs, err := client.Subscriptions.GetSubscriptionAssets(params)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	fmt.Println("Assets for subscription: " + cfg.subName)
// 	fmt.Println("")

// 	if len(subs) == 0 {
// 		fmt.Println("no subscription assets found")
// 		os.Exit(0)
// 	}

// 	for _, sub := range subs {
// 		fmt.Println(sub.Name + ", " + sub.License.Name)

// 		data := sub.Space

// 		eb, err := bytesize.Parse(fmt.Sprintf("%.0f %s", data.Effective.Data, data.Effective.Unit))
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}

// 		cb, err := bytesize.Parse(fmt.Sprintf("%.0f %s", data.Capacity.Data, data.Capacity.Unit))
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}

// 		fmt.Printf("Effective used: %s\n", eb)
// 		fmt.Printf("Used ratio: %.0f %s\n", data.UsedRatio.Data, data.UsedRatio.Unit)
// 		fmt.Printf("Capacity: %s\n", cb)
// 		fmt.Printf("Data reduction: %.0f %s\n", data.DataReduction.Data, data.DataReduction.Unit)

// 		fmt.Println("")
// 	}
// }
