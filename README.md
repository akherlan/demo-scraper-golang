# Ivosights Technical Test - News Scraper

## Overview

This project was developed as part of the technical test for my application to join Ivosights and was not production ready.

The objective of this project was to develop two news scrapers specifically designed for Detik.com and Liputan6.com, with the purpose of retrieving the most recent news articles available on these platforms.

## Accomplishments

I am pleased to report that I have successfully completed the homework test as requested. The project demonstrates my ability to:

1. Develop web scraper or crawler using Golang (which is a new experience for me and it's fun!)
2. Extract relevant information from news websites: title, content text, date published, link
3. Handle multiple sources: Detik.com and Liputan6.com
4. Architect and design scraping application with an easy-to-maintain structure
5. Utilize YAML for parsing configuration and management for flexibility and maintainability
6. Implement data persistence using MongoDB for storage

While I am pleased with the progress made, I acknowledge that there's always room for improvement and learning.

This program represents my sincere effort to apply newly acquired skills in Golang and I am grateful for the opportunity to demonstrate my adaptability and willingness to learn.

I hope that my work showcases potential value I could bring to the team at Ivosights.

## Technical Details

- **Programming language**: Golang
- **Target websites**:
  - Detik.com
  - Liputan6.com
- **Functionality**:
  - Collect the latest available news articles from both websites
  - Automatically creates unique IDs from article URLs to prevent duplicate entries in the database, safe for multiple times run
  - Provides clean news text with advertisement removal text content
- **Not implemented**:
  - JavaScript dynamic rendering or headless browser automation (to load more data from Liputan6 website)
  - Multiple pages crawling (focus on the latest news articles only and avoid to crawl all pages)
  - Manage scheduling for scraper
  - Options and arguments for advanced application control
  - Proxy settings

While certain advanced features were not implemented in this version, this decision was made to prioritize lightweight performance and simplicity for the purposes of technical test and local execution. It will keep the application simple and not requiring to use complex configuration to test and run.

However, the project was designed with future extensibility in mind. It hopfully can accommodate additional capabilities as needed, ensuring scalability and adaptability to evolving requirements.

## Project Structure

The project is structured to ensure modularity and improve organization to ease of maintenance and scalability.

```
scraper
├── config
│   ├── config.go       # loading and parsing config
│   └── config.yaml     # config file for parser, db, etc.
├── data.json           # sample data exported from db
├── db
│   └── mongo.go        # implement db connection and operation
├── go.mod              # go modules
├── go.sum
├── LICENSE.txt         # license file
├── main.go             # entry point of the app
├── news
│   └── article.go      # defines data struct, logic, model, method
├── README.md           # project explanation and documentation
└── sample.env          # sample environment variable
```

## Data Structure

```json
[
  {
    "_id": {
      "$oid": "6701efd4843fe9dc485bba06"
    },
    "content": "Jakarta - Massa dari Aliansi Rakyat Indonesia Bela Palestina (ARI-BP) menggelar aksi bela Palestina di depan Kedutaan Besar Amerika Serikat (Kedubes AS). Massa menyerukan penolakan terhadap standar ganda yang diberikan terhadap perang yang berkecamuk antara Israel dan Hamas di Jalur Gaza.\"Dan juga kita sama-sama berdiri di depan sebuah bangunan Kedutaan Besar atau boleh saja kita sebut Kedutaan Besar standar ganda dunia Amerika Serikat, kita di sini menolak segala bentuk standar ganda karena apa bedanya kita dengan saudara-saudara kita di Palestina. Ini bukan isu agama, gender ini adalah isu kemanusiaan,\" kata koordinator orasi dari Aliansi Pemuda Indonesia untuk Palestina, Abdullah Mudarik di depan Kedutaan Besar Amerika Serikat, Jakarta Pusat, Minggu (6/10/2024).Abdullah mengatakan aksi ini merupakan bentuk pengamalan terhadap UUD 1945. Dia mengatakan aksi ini juga bentuk cinta untuk negara dan warga Palestina. \"Maka kita berdiri di sini adalah bentuk cinta kita kepada negara kita dan bentuk cinta kepada saudara kita di Palestina,\" ujarnya Dia mengatakan terjadi ironi lantaran warga Palestina belum merasakan kemerdekaan. Dia mengajak massa melanjutkan perjuangan untuk Palestina tak berhenti usai aksi ini selesai digelar.\"Sungguh ironi kita bisa berteriak merdeka dengan lantang, tapi saudara kita di Palestina belum dan akan mencicipinya,\" ujar Abdullah.\"Mari kita berjanji kepada diri masing-masing bahwasanya dalam peringatan satu tahun ini perjuangan akan terus berlangsung,\" tambahnya.Dia juga mengajak massa menyerukan free Palestine. Massa mengikuti seruan tersebut.\"Free free Palestine, free free Palestine. From the river to the sea, Palestine will be free. Palestina Palestina, bebaskan bebaskan. Israel Israel, go to hell,\" ujarnya dan diikuti massa.Aksi itu mengusung tema 'Perjuangan Bersama Memperingati 1 Tahun Genosida di Gaza dan 76 Tahun Perlawanan Palestina'. Tak semua massa melakukan long march dari Monas.Sebagian massa sudah berkumpul di depan Kedubes AS. Sementara itu, Jalan Medan Merdeka Selatan dari dua arah yang mengarah ke Monas maupun ke Stasiun Gambir juga ditutup.Massa tampak mengenakan baju bernuansa putih dan hitam. Mereka juga membawa bendera Palestina dan poster bergambar ibu yang menggendong anaknya.Massa juga membawa poster bertulisan 'Stop Genosida'. Selain itu, massa juga membawa sorban dan ikat kepala bertuliskan 'Palestina'.Aksi dimulai dengan pembacaan doa dan dzikir. Mereka mendoakan warga Palestina.\"Baik, mari kita mulai, dzikir dan doa yang akan kita panjatkan. Kita mohon kepada Allah agar memberikan keberkahan, keridoan,\" kata Sekertaris Umum FPI, Buya Husein yang memimpin dzikir.Kemudian, massa juga menyanyikan lagu Indonesia Raya. Lalu, massa melantunkan ayat suci Alquran. (mib/dek) aksi bela palestina palestina kedubes as",
    "published": {
      "$date": "2024-10-06T02:03:00.000Z"
    },
    "title": "Massa Aksi Bela Palestina di Depan Kedubes AS Serukan Tolak Standar Ganda",
    "url": "https://news.detik.com/berita/d-7574171/massa-aksi-bela-palestina-di-depan-kedubes-as-serukan-tolak-standar-ganda"
  }
]
```

## Application Usage

The application requires Go and MongoDB installed on the operating system. Configure MongoDB connection with copying or renaming `sample.env` to `.env`, then edit to provide MongoDB connection string. Configuration for scraper settings provided on `./config/config.yaml` or `./config.yaml`.

Go module installation:

```shell
go mod download
```

The application can be run using one of following methods. First, directly run the `main.go` file:

```shell
go run main.go
```

Or second, compile and run the binary:

```shell
go build
./scraper
```

The scraper will collect articles from news websites and save results to MongoDB.

## License

Until discussed further, the project is currently under the MIT License. This licensing choice may be subject to review in the future as needed, particularly if the project's scope or usage evolves.

This project also utilizes various third-party modules and libraries, each subject to its own licensing terms. Please refer to the respective licenses of these dependencies for more information on their usage and distribution terms.

## Conclusion

The project showcases my technical skills and problem-solving abilities, which I hope align well with Ivosights' requirements. I am excited about the possibility of bringing these skills to your team and contributing to data-driven business at Ivosights.

Thank you for considering my application. I look forward to hear from you again.

-- [Andi Herlan](https://github.com/akherlan)

