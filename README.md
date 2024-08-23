### Cara ribet
- [Install golang](https://go.dev/doc/install)
- Unduh repo
```
$ git clone https://github.com/frmdeveloper/readsw
$ cd readsw 
```
- Install Module
```
$ go mod init frm
$ go mod tidy
```
- Ubah nomor Anda [(lokasinya)](https://github.com/frmdeveloper/readsw/blob/7d912976d52f984999b189adf72cff2f421a6e58/main.go#L74)<br>
kalo males, liat perintah `jalankan` di [cara simpel](https://github.com/frmdeveloper/readsw?tab=readme-ov-file#cara-simpel)
- Kemas dan jalankan
```
$ go build -o sw
$ ./sw
```
Kalo males ngemas
```
go run main.go
```

### Cara simpel
Download dan ekstrak binarynya [disini](https://github.com/frmdeveloper/readsw/releases/tag/1)
- jalankan
```
./sw -n 62831xxxxxxxx
```
