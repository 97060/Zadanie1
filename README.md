---
# Sprawozdanie - zadanie 1
### Autor: Michał Grabowiec 
---
## CZĘŚĆ OBOWIĄZKOWA

### 1. Kod serwera

Serwer został napisany w języku `Go`.

```go
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
	"time"
)
// Struktura używana do parsowania pliku JSON zwracanego przez API 
type Response struct {
	Meta struct {
		Code          string `json:"code"`
		ExecutionTime string `json:"execution_time"`
	} `json:"meta"`
	Data struct {
		IP          string      `json:"ip"`
		City        interface{} `json:"city"`
		Postal      interface{} `json:"postal"`
		State       interface{} `json:"state"`
		StateCode   interface{} `json:"state_code"`
		Country     string      `json:"country"`
		CountryCode string      `json:"country_code"`
		Location    string      `json:"location"`
		Timezone    struct {
			ID                  string `json:"id"`
			Location            string `json:"location"`
			CountryCode         string `json:"country_code"`
			CountryName         string `json:"country_name"`
			Iso31661Alpha2      string `json:"iso3166_1_alpha_2"`
			Iso31661Alpha3      string `json:"iso3166_1_alpha_3"`
			UnM49Code           string `json:"un_m49_code"`
			Itu                 string `json:"itu"`
			Marc                string `json:"marc"`
			Wmo                 string `json:"wmo"`
			Ds                  string `json:"ds"`
			PhonePrefix         string `json:"phone_prefix"`
			Fifa                string `json:"fifa"`
			Fips                string `json:"fips"`
			Gual                string `json:"gual"`
			Ioc                 string `json:"ioc"`
			CurrencyAlphaCode   string `json:"currency_alpha_code"`
			CurrencyCountryName string `json:"currency_country_name"`
			CurrencyMinorUnit   string `json:"currency_minor_unit"`
			CurrencyName        string `json:"currency_name"`
			CurrencyCode        string `json:"currency_code"`
			Independent         string `json:"independent"`
			Capital             string `json:"capital"`
			Continent           string `json:"continent"`
			Tld                 string `json:"tld"`
			Languages           string `json:"languages"`
			GeonameID           string `json:"geoname_id"`
			Edgar               string `json:"edgar"`
		} `json:"timezone"`
		Datetime struct {
			Date          string `json:"date"`
			DateTime      string `json:"date_time"`
			DateTimeTxt   string `json:"date_time_txt"`
			DateTimeWti   string `json:"date_time_wti"`
			DateTimeYmd   string `json:"date_time_ymd"`
			Time          string `json:"time"`
			Month         string `json:"month"`
			MonthWilz     string `json:"month_wilz"`
			MonthAbbr     string `json:"month_abbr"`
			MonthFull     string `json:"month_full"`
			MonthDays     string `json:"month_days"`
			Day           string `json:"day"`
			DayWilz       string `json:"day_wilz"`
			DayAbbr       string `json:"day_abbr"`
			DayFull       string `json:"day_full"`
			Year          string `json:"year"`
			YearAbbr      string `json:"year_abbr"`
			Hour12Wolz    string `json:"hour_12_wolz"`
			Hour12Wilz    string `json:"hour_12_wilz"`
			Hour24Wolz    string `json:"hour_24_wolz"`
			Hour24Wilz    string `json:"hour_24_wilz"`
			HourAmPm      string `json:"hour_am_pm"`
			Minutes       string `json:"minutes"`
			Seconds       string `json:"seconds"`
			Week          string `json:"week"`
			OffsetSeconds string `json:"offset_seconds"`
			OffsetMinutes string `json:"offset_minutes"`
			OffsetHours   string `json:"offset_hours"`
			OffsetGmt     string `json:"offset_gmt"`
			OffsetTzid    string `json:"offset_tzid"`
			OffsetTzab    string `json:"offset_tzab"`
			OffsetTzfull  string `json:"offset_tzfull"`
			TzString      string `json:"tz_string"`
			Dst           string `json:"dst"`
			DstObserves   string `json:"dst_observes"`
			TimedaySpe    string `json:"timeday_spe"`
			TimedayGen    string `json:"timeday_gen"`
		} `json:"datetime"`
	} `json:"data"`
}

const LOG_FILE = "./app.log" // Ścieżka do pliku z logami
const PORT = "8082" // Port na jakim działa serwer
const testIPAddress = "66.220.144.0" // Adres IP do celów pokazania działania serwera

func main() {
	
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE,0644)
	if err != nil {
        log.Panic(err)
    }
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// Wypisanie danych startowych do logu
	log.Println("Data uruchomienia  " + time.Now().Format("2006.01.02 15:04:05") + " Autor: Michał Grabowiec" + " PORT: " + PORT)

	// Obsługa odpowiednich ścieżek
	http.HandleFunc("/", serverHandler)
	http.HandleFunc("/log", logShowHandler)
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
}

// Główna funkcja zapewniająca funkcjonalność serwera 
func serverHandler(w http.ResponseWriter, r *http.Request) {
	msg := "" // Dane do wypisania
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	
	tmpStr := strings.Split(IPAddress, ":")
	IPAddress = tmpStr[0]
	log.Println("Nowe połączenie: " + IPAddress)
	if net.ParseIP(IPAddress).IsPrivate() || net.ParseIP(IPAddress).IsLoopback(){
		msg += "Ustawiono adres testowy  "
		IPAddress = testIPAddress
	}
	msg += " Adres IP: " + IPAddress
	url := "https://timezoneapi.io/api/ip/?" + IPAddress + "&token=abQGXSYEczrgJCKoQmuA" // API do pobrania daty i czasu dla danego IP
	response, err := http.Get(url)
	if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }
	responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
	var result Response
	json.Unmarshal(responseData, &result)
	msg += "\nData i godzina: " + result.Data.Datetime.Date + " " + result.Data.Datetime.Time
	fmt.Fprint(w, msg)
	
}

// Funkcja obsługująca wyświetlanie logów
func logShowHandler(w http.ResponseWriter, r *http.Request){
	content, err := ioutil.ReadFile("app.log")
    if err != nil {
        log.Fatal(err)
    }
	fmt.Fprint(w, string(content))
}
```
---

### 2. Dockerfile
```dockerfile
FROM golang:1.18 as gobuilder
WORKDIR /app
COPY server.go ./
COPY go.mod ./
COPY setup.sh ./
COPY upx-3.96-amd64_linux.tar.xz ./
RUN bash setup.sh && \
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server && \
apt-get update && \
apt-get install xz-utils && \
tar -C /usr/local -xf upx-3.96-amd64_linux.tar.xz && \
/usr/local/upx-3.96-amd64_linux/upx --ultra-brute --overlay=strip ./server

FROM scratch as main
LABEL Autor: "Michał Grabowiec"
COPY --from=gobuilder /app/server /
ADD ca-certificates.crt /etc/ssl/certs/
EXPOSE 8082
ENTRYPOINT [ "/server" ]
```
Plik [Dockerfile](../main/Dockerfile) wykorzystuje wieloetapową metodę budowania obrazu. Pierwsza warstwa odpowiada za zbudowanie pliku wykonywalnego serwera. Wykorzystywany jest pomocniczy skrypt znajdujący się w pliku [setup.sh](../main/setup.sh.sh), który umożliwia określenie architektury w celu budowania obrazów na wiele architektur. W celu zmiejszenia wielkości pliku wykonywalnego a co za tym idzie całego obrazu wykorzystywane są 2 metody:
- Przy kompilacji używane są flagi -w oraz -s, które usuwają informacje używane do debugowania oraz tablice symboli
- Po skompilowaniu używany jest packer [UPX](https://upx.github.io/), który dodatkowo zmniejsza wielkośc skompilowanego wcześniej pliku. Automatycznie wykrywa on architekture

Użyte metody pozwalają na redukcje obrazu do wielkości `1.68 MB`.
Druga warstwa jest wartswą scratch. Aby uruchomić serwer nie są wymagane dodatkowe zależności dlatego użycie tego typu warsty jest możliwe. Wymagane było jedynie dodanie [certyfikatów](../main/ca-certificates.srt) aby umożliwić komunikacje serwera z API. Serwer uruchamiany jest na porcie `8082`.  

---

### 3. Polecenia
**a.&ensp;Zbudowanie opracowanego obrazu kontenera:** 
```
DOCKER_BUILDKIT=1 docker build -t server . 
```

**b.&ensp; Uruchomienie kontenera ze zbudowanym obrazem**

```
docker run -t --name server -p 8082:8082 server
```

**c.&ensp; Działanie serwera i dostęp do danych**

Wchodząc w przeglądarkę i wpisując adres `localhost:8082` ukazuje nam się działająca strona uruchomionego serwera

![server](https://user-images.githubusercontent.com/103126350/167948226-ed0da1ad-ad28-4a17-b57c-836b3e4d2ea2.PNG)

Aby uzyskać dostęp do logów zapisywanych przez serwer należy użyć adresu `localhost:8082/log`

![logs](https://user-images.githubusercontent.com/103126350/167948265-638d5a9e-b91f-4c52-b75e-303ce0fc72de.PNG)


**d.&ensp; Sprawdzenie ilości warstw w zbudowanym obrazie**

```
docker image history server
```

![history](https://user-images.githubusercontent.com/103126350/167954330-42ddcfa4-2b55-4af6-9171-db2184fd88ca.PNG)


**Alternatywny sposób sprawdzenia ilości warstw oraz uzyskania innych informacji na temat zbudowanego obrazu**

```
docker image inspect server
```
![inspect](https://user-images.githubusercontent.com/103126350/167954338-da68f6d3-2e29-47ae-af8a-b7f99f06357d.PNG)


---

### 4. Budowanie obrazów na różne architektury

Aby było możliwe zbudowanie obrazów na różne platformy sprzętowe musimy skorzystać z zasobów emulatora `QEMU`. Na potrzeby wykonania tego zadania zainstalujemy `QEMU` lokalnie, ale można to zrobić w alternatywny sposób z wykorzystaniem dedykowanego kontenera. Następnie do zbudowania obrazów wykorzystamy wraper `buildx`. 

**Instalacja zasobów QEMU**

```
sudo apt-get install qemu-user-static
```

**Utworzenie nowego buildera buildx oraz ustawienie go jako domyślnego**

```
docker buildx create --name builder
```

```
docker buildx use builder
```

**Zbudowanie obrazu serwera na 3 wybrane platformy i przesłanie ich na repozytorium DockerHub**

```
docker buildx build -t 97060/zadanie1:multiplatform --platform linux/amd64,linux/arm64/v8,linux/arm/v7 --push . 
```
Z jakiegoś powodu występowały błędy podczas budowania na moim laptopie dlatego do zbudowania obrazów i dodania ich na DockerHub użyłem Github Actions.
```yml
- name: Build and push
        id: docker_build
        uses: docker/build-push-action@v2
        with:
          context: ./
          platforms: linux/amd64,linux/arm64/v8,linux/arm/v7
          file: ./Dockerfile
          push: true
          tags: |
            97060/zadanie1:multiplatform
            ghcr.io/97060/zadanie1:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
```
**Potwierdzenie poprawnego zbudowania obrazów - repozytorium `DockerHub`**

![DockerHub](https://user-images.githubusercontent.com/103126350/167954798-51ff9c02-9c42-46e2-9160-0c5d863756e0.PNG)

**Uruchomienie obrazu z repozytorium DockerHub**
```
docker run -p 8082:8082 97060/zadanie1:multiplatform
```


Zbudowane obrazy można znaleźć na repozytorium [DockerHub](https://hub.docker.com/r/97060/zadanie1/tags).

---
## CZĘŚĆ DODATKOWA

### **DODATEK 1**

### 1. GitHub Actions, Cache oraz Github Container Registry

Proces budowania może zostać zautomatyzowany przy pomocy `Github Actions`. Narzędzie to znajduje się na górnym panelu naszego repozytorium w zakładce `Actions`.

![githubActions](https://user-images.githubusercontent.com/103126350/167950021-efa29c50-f298-47b1-ac4a-40c6ae516064.PNG)

Nowy `workflow` dla Github Actions utowrzyłem dodając plik workflow.yml do katalogu .github w repozytoprium. 

![workflow](https://user-images.githubusercontent.com/103126350/167950201-078123a9-2a6b-48e7-ab14-36e5d687c4eb.PNG)


---

Na podstawie dokumentacji oraz przepływów tworzonych na laboratoriach utworzyłem [workflow](../.github/workflows/workflow.yml) budujący obraz na 3 zadanie architektury a następnie publikujący go na `GitHub Containers Registry` oraz `DockerHub`.
```yml
name: GitHub Actions workflow with push to GHCR.io and DockerHub

on:
  push:
    branches:
      - main
    paths-ignore:
      - '**/README.md'
      
jobs:
  build-push-images:
    name: Build and push to GHCR
    runs-on: ubuntu-latest
    
    steps:
        # sprawdzenie poprawności kodu
      - name: Checkout code
        uses: actions/checkout@v2

        # Uruchomienie QEMU
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v1

        # Uruchomienie buildx
      - name: Buildx set-up
        id: buildx
        uses: docker/setup-buildx-action@v1
        
        # Logowanie do GitHuba
      - name: Login to GitHub
        uses: docker/login-action@v1 
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GHCR_PASSWORD }}

        # Logowanie do DockerHub
      - name: Login to DockerHub
        uses: docker/login-action@v1 
        with:
          username: ${{ secrets.DOCKER_HUB_USERNAME }}
          password: ${{ secrets.DOCKER_HUB_ACCESS_TOKEN }}
    
        # Budowanie obrazów i ich publikacja
      - name: Build and push
        id: docker_build
        uses: docker/build-push-action@v2
        with:
          context: ./
          platforms: linux/amd64,linux/arm64/v8,linux/arm/v7
          file: ./Dockerfile
          push: true
          tags: |
            97060/zadanie1:multiplatform
            ghcr.io/97060/zadanie1:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max


```    
Po każdorazowym wykonaniu operacji `git push` uruchamiany jest stworzony przepływ.

![start](https://user-images.githubusercontent.com/103126350/167973778-199eeb2f-0727-48fe-bd1e-5681ecb8cf51.PNG)


Otrzymujemy informacje czy przepyw zakoczył się prawidłowo czy wystąpił błąd. 

![successError](https://user-images.githubusercontent.com/103126350/167955031-43b883f4-fd5e-43ca-8200-616b5f77b482.PNG)

Wchodząc na nasz profil `GitHub`, a następnie w zakładkę `Packages` widzimy paczkę zawierającą zbudowane obrazy.

![package](https://user-images.githubusercontent.com/103126350/167955204-29ffbfc1-8006-420d-9b5b-d42ef4cd1518.PNG)



---

### Użycie cache

Na poniższym zrzucie ekranu doskonale widać różnice między budowaniem obrazu z użyciem cache i bez. W tym przypadku różnica wynosi prawie 4 minuty.

![cache](https://user-images.githubusercontent.com/103126350/167968818-aa8f161e-9ba7-47f7-9f18-974dc197721a.PNG)

Przy budowaniu kiedy nie ma jeszcze cache jest wykonywany export.

![export](https://user-images.githubusercontent.com/103126350/167968892-f4304bee-a9b9-4c1e-9217-3ac5deb2c3a3.PNG)

Przy budowaniu kiedy cache już istnieje warstwy, które zostały już utworzone i są w pamięci cache nie są tworzone ponownie.

![cached](https://user-images.githubusercontent.com/103126350/167969005-ccad12bd-1d5e-4913-bd7b-04f104131926.PNG)



### **DODATEK 2**

### 1. Uruchomienie prywatnego rejestru 

**a.&ensp;Uruchomienie kontenera na porcie 6677** 

Uruchomienie kontenera z rejestrem. Flaga `--restart=always` powoduje, że kontener w razie zatrzymania automatycznie zostanie zrestartowany i uruchomiony ponownie. 
```
docker run -d -p 6677:5000 --restart=always --name private_registry registry
```

**b.&ensp;Pobranie najnowszego Ubuntu i wgranie go do utworzonego rejestru** 

```
docker pull ubuntu:latest
```

Zmiana nazwy obrazu - dodanie odpowiedniego tagu, który oznacza równiez lokalizacje dla obrazu przy pushowaniu.

```
docker tag ubuntu:latest localhost:6677/ubunciak_registry
```

Wgrywanie obrazu do rejestru.

```
docker push localhost:6677/ubunciak_registry
```

Dodany obraz jest widoczny w rejestrze i może być z niego pobrany.

![reg](https://user-images.githubusercontent.com/103126350/167959830-e26a6d54-4a29-417c-a1da-e202a9b5e9f9.PNG)

![ubunciak](https://user-images.githubusercontent.com/103126350/167972408-870d2d18-5f17-4995-9ceb-fc009aa42cfb.PNG)



---

### 2. Dodanie mechanizmu kontroli dostępu htpasswd. 

Aby korzystać z uwierzytelniania należy wcześniej skonfigurować `TLS`co wiąże się z wygenerowaniem odpowiedniego certyfikatu dla `localhost`.

Tworzymy foldery `cert` oraz `auth`

```
mkdir cert
mkdir auth
```

W folderze certs generujemy odpowiednie klucze
```
cd cert/
openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout cert.key -out cert.pem
openssl x509 -outform pem -in cert.pem -out cert.crt
openssl req -new -nodes -newkey rsa:2048 -keyout localhost.key -out localhost.csr
```

Tworzymy plik `domains.ext` z podaną zawartością:
```
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
```

Generujemy plik locahost.crt
```
openssl x509 -req -sha256 -days 1024 -in localhost.csr -CA cert.pem -CAkey cert.key -CAcreateserial -extfile domains.ext -out localhost.crt
```
Generujemy dane do autentykacji w folderze auth
```
cd ..
docker run --entrypoint htpasswd httpd:2 -Bbn admin admin > auth/htpasswd
```
Wyłączamy poprzednio uruchomiony rejestr
```
docker container stop private_registry
docker rm private_registry
```

Uruchamiamy kontener z użyciem uwierzytelnienia:
```
docker run -d \
   -p 6677:5000 \
   --restart=always \
   --name private_registry \
   -v "$(pwd)"/auth:/auth \
   -e "REGISTRY_AUTH=htpasswd" \
   -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
   -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
   -v "$(pwd)"/cert:/cert \
   -e REGISTRY_HTTP_TLS_CERTIFICATE=/cert/localhost.crt \
   -e REGISTRY_HTTP_TLS_KEY=/cert/localhost.key \
   registry
```
Teraz dodanie obrazu do lokalnego reporytorium jest możliwe jedynie po uwierzytelnieniu loginem i hasłem odpowiedniego użytkownika:

![docker logowanie](https://user-images.githubusercontent.com/103126350/167971693-fb2143a9-e966-46ff-baa0-81ab3115f17b.PNG)

---
