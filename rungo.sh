export GOROOT=/usr/local/go
export GOPATH=/root/gopath
export GOBIN=/root/gopath/bin

rm -fr hdwallet
go build -o hdwallet ./
./hdwallet -C ./hdwallet.conf
# Your wallet generation seed is:
# 5a2fbfca86818e418a71646627da53f54383b8bf3e330a97926f9e714f078659
# IMPORTANT: Keep the seed in a safe place as you
# will NOT be able to restore your wallet without it.
# Please keep in mind that anyone who has access
# to the seed can also restore your wallet thereby
# giving them access to all your funds, so it is
# imperative that you keep it in a secure location.