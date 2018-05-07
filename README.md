# walletRecover #

walletRecover is a utility to help recover legacy blockchain.info wallets.
It supports multiple legacy blockchain.info encryption formats

## Installation ##

	go get -u github.com/solipsis/walletRecover
	go install github.com/solipsis/walletRecover
	
	Or if you'd prefer to not build from source, you can find precompiled binaries in the releases tab
  
## Usage ##

```
walletRecover [pathToEncryptedWallet] [pathToPasswordDictionary]

example: walletRecover myWallet.aes.json rockyou.txt
```

### Example Output ###

```
Health Check: Not json: P5a4t1r9i0c9k5 | 0 passwords tried so far
Health Check: Not json: 123456n | 10000 passwords tried so far
Health Check: Not json: ineedu | 20000 passwords tried so far
Health Check: Not json: angies | 30000 passwords tried so far
Health Check: Not json: nathan4 | 40000 passwords tried so far
Health Check: Not json: alphonso | 50000 passwords tried so far
Health Check: Not json: monmouth | 60000 passwords tried so far
Health Check: Not json: shawarma | 70000 passwords tried so far
Health Check: Not json: 300606 | 80000 passwords tried so far
Health Check: Not json: alethea | 90000 passwords tried so far
Health Check: Not json: salah | 100000 passwords tried so far
Health Check: Not json: ronaldinio | 110000 passwords tried so far
Health Check: Not json: 30011991 | 120000 passwords tried so far
Health Check: Not json: num1bitch | 130000 passwords tried so far
Health Check: Not json: holton | 140000 passwords tried so far
Health Check: Not json: dunkel | 150000 passwords tried so far
Health Check: Not json: fucku07 | 160000 passwords tried so far
Health Check: Not json: lane05 | 170000 passwords tried so far
Health Check: Not json: sexygirl8 | 180000 passwords tried so far
Health Check: Not json: Auburn | 190000 passwords tried so far
Health Check: Not json: juliepearl | 200000 passwords tried so far
Health Check: Not json: 05142005 | 210000 passwords tried so far
Health Check: Not json: ilikepie12 | 220000 passwords tried so far
Health Check: Not json: 12140 | 230000 passwords tried so far
Health Check: Not json: lokita5 | 240000 passwords tried so far
Health Check: Not json: ambriz | 250000 passwords tried so far
Health Check: Not json: snowfire | 260000 passwords tried so far
Health Check: Not json: joecole10 | 270000 passwords tried so far
Health Check: Not json: alovetokill | 280000 passwords tried so far
Wallet decoded successfully with password: "potato"
Decoded: {
	"guid" : "695b1421-b294-10a2-8456-15f4508f9f14",
	"sharedKey" : "084a5323-2fe2-4e1e-8d91-8edacb7aa957",
	"options" : {"pbkdf2_iterations":10,"fee_policy":0,"html5_notifications":false,"logout_time":600000,"tx_display":0,"always_keep_local_backup":false},
	"keys" : [
	{"addr" : "1NbKx2bfD3KXYyacbNWeZy7UAeBBqskFic",
	 "priv" : "BuripA8RrhFnW2WMsHiqEEE4itW35gtSJz7d24J4SFUk"}
	]
}
```

### Donations ###
```
BTC: 1ED8d3r93XNRs5za1xxeNcTp8ugeoZfYUx
ETH: 0xdd97359867fda3713def38eeb9390994de206be2
```
