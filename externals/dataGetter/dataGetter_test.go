package dataGetter

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDataFromFiles(t *testing.T) {
	baseDir := "/tmp/blockchian-data-aggregator-datas"
	content := `"app","ts","event","project_id","source","ident","user_id","session_id","country","device_type","device_os","device_os_ver","device_browser","device_browser_ver","props","nums"
"seq-market","2024-04-15 02:15:07.167","BUY_ITEMS","4974","","1","0896ae95dcaeee38e83fa5c43bef99780d7b2be23bcab36214","5d8afd8fec2fbf3e","DE","desktop","linux","x86_64","chrome","122.0.0.0","{""tokenId"":""215"",""txnHash"":""0xd919290e80df271e77d1cbca61f350d2727531e0334266671ec20d626b2104a2"",""chainId"":""137"",""collectionAddress"":""0x22d5f9b75c524fec1d6619787e582644cd4d7422"",""currencyAddress"":""0xd1f9c58e33933a993a3891f8acfe05a68e1afc05"",""currencySymbol"":""SFL"",""marketplaceType"":""amm"",""requestId"":""""}","{""currencyValueDecimal"":""0.6136203411678249"",""currencyValueRaw"":""613620341167824900""}"
"seq-market","2024-04-15 02:26:37.134","BUY_ITEMS","4974","","1","0896ae95dcaeee38e83fa5c43bef99780d7b2be23bcab36214","5d8afd8fec2fbf3e","DE","desktop","linux","x86_64","chrome","122.0.0.0","{""currencyAddress"":""0xd1f9c58e33933a993a3891f8acfe05a68e1afc05"",""currencySymbol"":""SFL"",""marketplaceType"":""amm"",""requestId"":"""",""tokenId"":""602"",""txnHash"":""0x1133d2837267e0de2eddf3655a3df99e055d172cb53c4e8e108e70322438e994"",""chainId"":""137"",""collectionAddress"":""0x22d5f9b75c524fec1d6619787e582644cd4d7422""}","{""currencyValueDecimal"":""2.361412166673735"",""currencyValueRaw"":""2361412166673735000""}"
"seq-market","2024-04-15 02:42:32.507","BUY_ITEMS","4974","","1","0896ae95dcaeee38e83fa5c43bef99780d7b2be23bcab36214","73ffe889f5223b5e","DE","desktop","linux","x86_64","chrome","122.0.0.0","{""marketplaceType"":""amm"",""requestId"":"""",""tokenId"":""201"",""txnHash"":""0x6c51abf80365cbf6a8a03d9e5fe939712742dff4b088d4f99ba44551907e5c2f"",""chainId"":""137"",""collectionAddress"":""0x22d5f9b75c524fec1d6619787e582644cd4d7422"",""currencyAddress"":""0xd1f9c58e33933a993a3891f8acfe05a68e1afc05"",""currencySymbol"":""SFL""}","{""currencyValueRaw"":""364528625334421950"",""currencyValueDecimal"":""0.36452862533442193""}"
"seq-market","2024-04-15 11:24:12.561","BUY_ITEMS","4974","","1","0896ae95dcaeee38e83fa5c43bef99780d7b2be23bcab36214","7a5bfa068c342272","DE","desktop","linux","x86_64","chrome","122.0.0.0","{""tokenId"":""601"",""txnHash"":""0xcd5e34370546c26bc426bcfda6fcfc8fd0e08d2be6ee2e00f4d6d455318f8640"",""chainId"":""137"",""collectionAddress"":""0x22d5f9b75c524fec1d6619787e582644cd4d7422"",""currencyAddress"":""0xd1f9c58e33933a993a3891f8acfe05a68e1afc05"",""currencySymbol"":""SFL"",""marketplaceType"":""amm"",""requestId"":""""}","{""currencyValueDecimal"":""3.2776439715275094"",""currencyValueRaw"":""3277643971527509500""}"`

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Fatal(err)
	}

	tmpfile, err := os.CreateTemp(baseDir, "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var wg sync.WaitGroup
	g := New(baseDir, 1)
	wg.Add(2)

	count := 0
	go func(t *testing.T, g *DataGetter) {
		defer wg.Done()
		for {
			select {
			case _ = <-g.Channel():
				count++
			case _ = <-g.EndChannel():
				assert.Equal(t, 4, count)
				return
			}
		}
	}(t, g)

	go func() {
		defer wg.Done()
		if err := g.ReadDataFromFiles(); err != nil {
			t.Errorf("failed to read data from files: %v", err)
		}
	}()
	wg.Wait()
}
