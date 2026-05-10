# PROPOSAL PROYEK AKHIR

# RANCANG BANGUN SISTEM PRESENSI FACE

# RECOGNITION BERBASIS ALGORITMA ROBUST PCA DAN

# GA-SVM PADA RASPBERRY PI DENGAN INTEGRASI JIRA

# FARCHAN PUTRA INDRIANTO

## TEKNOLOGI REKAYASA KOMPUTER

## SEKOLAH VOKASI

## INSTITUT PERTANIAN BOGOR

## BOGOR

## 2025


# PERNYATAAN MENGENAI LAPORAN PROYEK AKHIR DAN

# SUMBER INFORMASI SERTA PELIMPAHAN HAK CIPTA

Dengan ini saya menyatakan bahwa Laporan Proyek Akhir dengan judul
“Rancang Bangun Sistem Presensi _Face Recognition_ Berbasis Algoritma _Robust_
PCA dan GA-SVM pada Raspberry Pi dengan Integrasi Jira” adalah karya saya
dengan arahan dari dosen pembimbing dan belum diajukan dalam bentuk apa pun
kepada perguruan tinggi mana pun. Sumber informasi yang berasal atau dikutip dari
karya yang diterbitkan maupun tidak diterbitkan dari penulis lain telah disebutkan
dalam teks dan dicantumkan dalam Daftar Pustaka di bagian akhir Laporan Proyek
Akhir ini.
Dengan ini saya melimpahkan hak cipta dari karya tulis saya kepada Institut
Pertanian Bogor.

```
Bogor, 12 November 2025
```
```
Farchan Putra Indrianto
J
```

```
dalam bentuk apa pun tanpa izin IPB.
```
```
Laporan Proyek Akhir
sebagai salah satu syarat untuk memperoleh gelar
Sarjana Terapan pada
Program Studi Teknologi Rekayasa Komputer
```
# Rancang Bangun Sistem Presensi Face Recognition Berbasis

# Algoritma Robust PCA dan GA-SVM pada Raspberry Pi

# dengan Integrasi Jira

# FARCHAN PUTRA INDRIANTO

## TEKNOLOGI REKAYASA KOMPUTER

## SEKOLAH VOKASI

## INSTITUT PERTANIAN BOGOR

## BOGOR

## 2025


Judul Proyek Akhir : Rancang Bangun Sistem Presensi _Face Recognition_
Berbasis Algoritma _Robust_ PCA dan GA-SVM pada
Raspberry Pi dengan Integrasi Jira
Nama : Farchan Putra Indrianto
NIM : J

```
Disetujui oleh
```
```
Pembimbing 1:
Prof. Dr. Ir. Sri Nurdiati, M.Sc.
196011261986012001
```
#### __________________

```
Diketahui oleh
```
```
Ketua Program Studi:
Dr. Inna Novianty, M.Si
201811198611192014
```
#### __________________

```
Tanggal Presentasi:
```
-


## DAFTAR ISI

DAFTAR TABEL iv


- I PENDAHULUAN DAFTAR GAMBAR iv
   - 1.1 Latar Belakang
   - 1.2 Rumusan Masalah
   - 1.3 Tujuan
   - 1.4 Manfaat
   - 1.5 Ruang Lingkup
- II TINJAUAN PUSTAKA
   - 2.1 Digitalisasi Manajemen Sumber Daya Manusia (SDM)
   - 2.2 Ergonomi Kognitif dalam Interaksi Sistem
   - 2.3 Sistem Pengenalan Wajah
   - 2.4 Algoritma Robust PCA dan GA-SVM
   - 2.5 Arsitektur Raspberry Pi
   - 2.6 Integrasi Sistem API Jira
- III METODE
   - 3.1 Lokasi dan Waktu Proyek Akhir
   - 3.2 Teknik Pengumpulan Data dan Analisis Data
   - 3.3 Prosedur Kerja
   - 3.4 Daftar Komponen
   - 3.5 Desain Sistem
   - 3.6 Matriks Rencana Proyek Akhir
- DAFTAR PUSTAKA
- 1 Spesifikasi dan fitur Raspberry Pi DAFTAR TABEL
- 2 Daftar komponen perangkat keras
- 3 Daftar komponen perangkat lunak
- 4 Jadwal pelaksanaan proyek akhir
- 1 Pinout Raspberry Pi DAFTAR GAMBAR
- 2 Prosedur kerja
- 3 Skematik rangkaian
- 4 Diagram flowchart fase pelatihan data
- 5 Diagram flowchart fase real-time
- 6 Bentuk 3D alat presensi
- 7 Peletakan kamera pada alat presensi


# I PENDAHULUAN

### 1.1 Latar Belakang

Manajemen sumber daya manusia pada era industri 5.0 merupakan suatu
bentuk pemanfaatan teknologi tingkat lanjut dalam upaya peningkatan mekanisme,
prosedur, hingga sistem secara keseluruhan. Keakuratan dan efisiensi dalam
mencatat kehadiran karyawan merupakan faktor kunci untuk mengoptimalkan
produktivitas dan efektivitas organisasi (Saied dan Syafii 2023). Salah satu dari
pemanfaatan tersebut adalah sistem presensi karyawan menggunakan teknologi
terbaru.
PT. Permodalan Nasional Madani merupakan lembaga BUMN (Badan Usaha
Milik Negara) sebagai pengelola keuangan dalam memajukan, memelihara maupun
mengembangkan Usaha Mikro, Kecil dan Menengah (UMKM). Sistem presensi
yang saat ini dipakai oleh PT. Permodalan Nasional Madani berbasis pindai kode
QR melalui aplikasinya, yaitu PNM Digi. Walaupun sudah mendigitalisasi, sistem
presensi berbasis pindah kode QR ini memiliki kelemahan signifikan dari sisi
ergonomi interaksi sistem. Proses ini menuntut beberapa langkah manual dari
karyawan, seperti mengharusnya mengeluarkan ponsel dan memposisikan kamera
untuk memindai. Alur kerja yang tidak _seamless_ dapat menyebabkan banyak friksi
yang dapat dinilai tidak efisien dan menurunkan kenyamanan karyawan.
Salah satu alur bisnis yang krusial di dalam internal PT. PNM adalah adanya
validitas KPI ( _Key Performance Indicator_ ). Saat ini, salah satu variabel KPI
tahunan, yaitu jumlah realisasi _task_ (tiket) masih dicatat secara manual oleh
karyawan. Proses manual ini memiliki dua kelemahan utama. Pertama, data tersebut
rentan terhadap ketidakakuratan, baik secara sengaja maupun tidak disengaja.
Kedua, proses verifikasi menjadi tidak efisien. Dengan struktur manajerial berlapis,
pemeriksaan manual untuk mencocokkan data laporan dengan data aktual di Jira
untuk setiap karyawan dinilai tidak efisien. Oleh karena itu, penelitian ini
mengusulkan sebuah sistem presensi terintegrasi yang dirancang untuk menjawab
kedua permasalahan tersebut.
Pengenalan wajah adalah salah satu aplikasi yang paling berkembang dalam
_computer vision_ dan _pattern recognition_ (Susim dan Darujati 2021). Pengenalan
wajah yang akan diimplementasikan pada penelitian ini berbasis algoritma _Robust_
PCA _(Principal Component Analysis)_ dan GA dengan integrasi _platform_ Jira. Fitur
dari citra wajah akan diekstraksi menggunakan algoritma _Robust PCA_ , yang akan
memisahkan komponen citra yang mengalami oklusi (terhalang) dan yang tidak
mengalami oklusi (Eman _et al_. 2023) melalui komputer _single-board_ Raspberry Pi.
Selanjutnya, fitur tersebut diklasifikasikan menggunakan _Support Vector Machine_
(SVM) yang parameternya dioptimalkan oleh algoritma genetika (GA) untuk
mencapai akurasi maksimal. Menurut Subiyanto _et al_. (2020), sistem pengenalan
wajah berbasis PCA-GA dapat meraih tingkat akurasi sebesar 90% dengan waktu
yang dibutuhkan untuk melakukan pelatihan selama 465,05 detik terhadap 116 citra
latih.
Secara bersamaan, untuk mengatasi masalah validitas data kinerja, sistem ini
akan diintegrasikan dengan API Jira. Integrasi ini dirancang untuk mengotomatisasi
alur perhitungan dan pelaporan data realisasi _task_ karyawan secara langsung saat
mereka melakukan presensi.


### 1.2 Rumusan Masalah

Implementasi sistem presensi pengenalan wajah berbasis _machine learning_
IoT memiliki beberapa parameter keberhasilan yang perlu dicapai. Parameter
tersebut mencakupi hal-hal internal dan eksternal. Perihal internal melibatkan
mesin dan algoritma yang digunakan untuk _machine learning_ , sedangkan perihal
eksternal melibatkan karyawan PT. PNM secara keseluruhan. Berikut merupakan
tiga rumusan masalah dari projek akhir ini:

1. Bagaimana merancang arsitektur _Cloud-Edge Computing_ pada sistem
    presensi pengenalan wajah yang efisien dalam penggunaan _bandwidth_ dan
    memiliki skalabilitas tinggi?
2. Berapa tingkat akurasi pengenalan wajah dan rata-rata waktu pemrosesan
    ( _latency_ ) yang dibutuhkan sistem secara _real-time_?
3. Bagaimana rancang bangun integrasi Jira dapat mengotomatisasi alur
    perhitungan dan pelaporan data realisasi _task_ di KPI ( _Key Performance_
    _Indicator_ ) karyawan PT. PNM?

### 1.3 Tujuan

Berdasarkan pendahuluan dan rumusan masalah yang telah diciptakan,
terdapat tiga tujuan yang ingin dicapai sebagai tingkat keberhasilan dari penelitian.
Berikut merupakan tiga tujuan pada hasil dari penelitian ini:

1. Merancang arsitektur Cloud-Edge Computing pada sistem presensi
    pengenalan wajah untuk memiliki penggunaan _bandwidth_ yang kecil dan
    skalabilitas tinggi.
2. Menganalisis waktu pemrosesan ( _latency_ ) sistem untuk implementasi
    secara _real-time_.
3. Mengintegrasikan Jira dengan sistem presensi dalam rangka
    mengotomatisasikan perhitungan dan pelaporan data realisasi _task_ di KPI
    ( _Key Performance Indicator_ ) karyawan PT. PNM.

### 1.4 Manfaat

Diharapkan bahwa projek akhir ini akan memberikan manfaat akademik
praktis maupun teoritis. Secara praktis, sistem yang dikembangkan akan
meningkatkan efisiensi dan tingkat ergonomi interaksi sistem presensi karyawan di
PT. PNM yang otomatis dan _seamless_. Dengan mengintegrasikan API Jira, sistem
ini mengotomatisasi alur perhitungan dan pelaporan data realisasi _task_. Integrasi ini
dapat menghilangkan proses laporan manual yang rentan tidak akurat dan sulit
diverifikasi, sehingga menyediakan data kinerja yang lebih akurat dan dapat diaudit
oleh manajemen.
Projek akhir ini akan menyongsong penggunaan teknologi tingkat lanjut
dalam sektor manajemen sumber daya manusia. Implementasi dan hasil dapat
digunakan sebagai referensi untuk merancang sistem serupa di bidang terkait yang
membutuhkan otomatisasi pemantauan harian.
Secara akademis, penelitian ini berkontribusi dalam bidang rekayasa sistem
cerdas untuk manajemen sumber daya manusia. Inovasi penelitian ini terletak pada
penerapan kombinasi algoritma _Robust_ PCA dan GA-SVM untuk menciptakan
sistem pengenalan wajah yang akurat. Selain itu, hasil penelitian ini diharapkan
dapat menjadi rujukan untuk memotivasi eksplorasi lebih lanjut mengenai integrasi


data presensi dengan alat ukur produktivitas kerja seperti Jira untuk menciptakan
lingkungan kerja yang lebih efisien dan terukur.

### 1.5 Ruang Lingkup

Penelitian ini berfokus pada rancang bangun dan implementasi sistem pada
_single-board computer_ Raspberry Pi 4 model B dengan jumlah RAM 4GB. Seluruh
pengujian, khususnya latensi pemrosesan, terikat pada _platform_ perangkat keras
tersebut. Dari sisi algoritma, penelitian ini terbatas pada penggunaan _Robust_ PCA
untuk ekstraksi fitur dan GA-SVM untuk klasifikasi. Dengan demikian, penelitian
ini tidak melakukan perbandingan performa dengan metode _Deep Learning_ (seperti
CNN) atau algoritma _machine learning_ klasik lainnya. Meskipun algoritma yang
lebih modern seperti _Gradient Boosting_ dan _Deep Learning_ (seperti CNN) tersedia,
SVM dipilih secara spesifik karena efisiensi inferensinya yang tinggi, yang krusial
untuk mencapai latensi pemrosesan _real-time_ pada _platform edge computing_
dengan sumber daya terbatas dalam Raspberry Pi 4.
Ruang lingkup fungsionalitas integrasi API Jira dibatasi hanya pada
pengambilan data ( _read-only_ ) untuk mengotomatisasi perhitungan data realisasi
_task_ , dan tidak mencakup modifikasi data di Jira atau perhitungan komponen KPI
bisnis lainnya. Kondisi operasional sistem diuji dalam lingkungan pencahayaan
dalam ruangan ( _indoor_ ) yang memadai dengan pose wajah menghadap ke depan
( _frontal_ ). Terakhir, _dataset_ yang digunakan untuk melatih dan menguji model
merupakan _dataset_ pribadi yang diambil di PT. PNM, bukan _dataset_ publik standar.


## II TINJAUAN PUSTAKA

### 2.1 Digitalisasi Manajemen Sumber Daya Manusia (SDM)

Keberhasilan suatu organisasi sangat dipengaruhi oleh kinerja individu
karyawannya. Setiap organisasi, termasuk sekolah, akan selalu berusaha untuk
meningkatkan kinerja karyawan, dengan harapan apa yang menjadi tujuan
perusahaan akan tercapai (Susilo dan Abdurrahman 2023). Kunci-kunci
keberhasilan tersebut dapat didorong dengan adanya digitalisasi pada manajemen
sumber daya manusia. Terutama apabila dorongan teknologi ini mempengaruhi
sistem presensi pada perusahaan atau industri.
Manajemen karyawan digital adalah tentang perencanaan dan penerapan
teknologi digital untuk mendukung dan membangun jaringan profesi SDM
(Setianingrum 2024). Peran teknologi tidak lagi sebatas menggantikan tugas
manual, melainkan untuk meningkatkan efisiensi, akurasi data, dan pengalaman
karyawan secara keseluruhan.
Transformasi ini memiliki pengaruh signifikan pada sistem presensi.
Berdasarkan studi yang dilakukan oleh Susilo dan Abdurrahman (2023) di
lingkungan kerja sebuah Madrasah Aaliyah (MA), presensi manual yang telah di
terapkan sebelumnya dinilai tidak efektif karena dapat pencatatannya dapat
dimanipulasi dan kurang terekam dengan baik. Hal ini menguatkan agenda untuk
mendigitalisasi infrastruktur manajemen karyawan. Selain itu, berdasarkan
penelitian Meutia _et al._ (2022), jika absensi setiap karyawan meningkat maka
kinerja karyawan tersebut juga akan meningkat.

### 2.2 Ergonomi Kognitif dalam Interaksi Sistem

Ergonomi memberi kemudahan untuk manusia dalam berbagai hal di dalam
lingkungan kerja, sehingga manusia memiliki kemudahan, kenyamanan serta
efisiensi dalam melakukan pekerjaan (Rahma dan Astuti 2025). Ergonomi sendiri
berfokus pada interaksi antara manusia dan elemen sistem atau lingkungan kerja,
dan bagaimana bisa dioptimalkan. Khususnya pada era digitalisasi dan
perkembangan teknologi terkini, ergonomi kognitif menghadapi tantangan baru.
Sistem berbasis _Internet of Things_ (IoT), kecerdasan buatan (AI), dan
otomatisasi semakin banyak digunakan, baik di lingkungan kerja maupun dalam
kehidupan sehari-hari (Bora _et al._ 2025). Menurut Waliyaden dan Leo (2024),
interaksi manusia dengan teknologi semakin kompleks, dan peran ergonomi
kognitif menjadi semakin penting untuk memastikan bahwa teknologi tersebut
mendukung, bukan menghambat, kemampuan manusia. Dalam peran digitalisasi
sistem presensi, khususnya sistem pindai kode QR. Proses yang mengharuskan
karyawan untuk melakukan upaya lebih, seperti membuka ponsel dan pindai kode
QR, dinilai kurang memenuhi kriteria psikologis (menambah beban kognitif) dan
kriteria hasil kerja (kurang efisien).
Oleh karena itu, penerapan sistem presensi yang terotomatisasi, seperti
pengenalan wajah, menjadi sebuah ideal secara prinsip ergonomi. Perubahan dari
proses aktif pindai kode QR menjadi pasif menampilkan wajah ke kamera
mengurangi upaya fisik dan mental, sehingga dapat meningkatkan kenyamanan
karyawan dan berdampak positif terhadap ergonomi kognitif karyawan.


### 2.3 Sistem Pengenalan Wajah

Sistem pengenalan wajah adalah teknologi biometrik yang menggunakan
wajah manusia untuk identifikasi dan verifikasi sesuatu (Razaqa _et al._ 2024). Wajah
manusia memiliki ragam keunikan dan karakteristik yang berbeda untuk setiap
individu. Keunikan karakteristik wajah setiap individu menjadi sebuah relevansi
ketika ingin membuat identifikasi diri sendiri dengan foto wajah maupun sidik jari.
Struktur sistem pengenalan wajah itu mirip seperti struktur sistem biometrik,
karena melibatkan deteksi wajah, _preprocessing_ citra wajah, ekstraksi fitur wajah,
dan klasifikasi fitur wajah (Oloyede _et al._ 2020). Oloyode _et al._ (2020) menjelaskan
lebih lanjut tahapan sistem pengenalan wajah:

1. Deteksi wajah merupakan proses verifikasi presensi wajah manusia.
2. _Preprocessing_ citra wajah adalah proses menyiapkan gambar agar berisi
    fitur wajah yang penting saja. Beberapa cara yang dilakukan adalah
    _normalization_ (transformasi skala citra wajah) dan _face alignment_
    (menentukan titik-titik fidusial atau _landmark_ pada citra wajah).
3. Ekstraksi fitur wajah adalah proses pengekstraksian data visual wajah
    paling relevan yang dapat mengidentifikasi wajah manusia secara unik.
4. Klasifikasi fitur wajah merupakan tahap pengenalan ( _recognition_ ) pada
    citra wajah, di mana sebuah citra wajah akan dicocokkan dengan data
    yang tersimpan di dalam _database_ untuk tujuan verifikasi atau identifikasi.
Menurut Mirghani dan Al-Mazruii (2022), sistem pengenalan wajah telah
diimplementasikan dalam berbagai aplikasi, menjadikannya sebuah teknologi
_computer vision_ yang krusial dalam era pengembangan teknologi, khususnya dalam
deteksi wajah dan klasifikasi fitur. Dalam penelitian ini, fokus utama akan diberikan
pada tahap ekstraksi fitur menggunakan algoritma _Robust_ PCA dan tahap
klasifikasi fitur yang dioptimalkan dengan GA-SVM.

### 2.4 Algoritma Robust PCA dan GA-SVM

Menurut Li _et al_. (2020), pengenalan wajah merupakan masalah pengenalan
pola visual dimana pemasukan visual direpresentasikan dalam bentuk matriks di
komputer perlu membedakan jika data tersebut memiliki sebuah wajah, lalu
mengidentifikasi siapa yang mempunyai wajah tersebut. Dalam proses ini, citra
wajah direpresentasikan sebagai matriks numerik yang dapat diolah oleh komputer.
Maka dari itu, peran algoritma sangat krusial dalam analisa dan identifikasi. Dalam
konteks pengerjaan penelitian ini, algoritma utama yang dipakai adalah _Robust_
PCA ( _Principal Component Analysis_ ) dan algoritma genetika berbasis SVM
( _Support Vector Machine_ ).

```
Algoritma Robust PCA ( Principal Component Analysis )
Principal Component Analysis (PCA) merupakan salah satu metode
berbasis penampilan yang popular digunakan untuk mereduksi dimensi dari
sekumpulan atau ruang citra sehingga basis atau koordinat yang baru dapat
menggambarkan model yang khas dari kumpulan tersebut (Sari et al. 202 3).
Pada computer vision , prinsip ini digunakan untuk representasi citra dengan
dimensi vektor fitur yang relatif rendah. Namun menurut Eman et al. (2023),
PCA sensitif terhadap pencilan, yang dapat mendistorsi hasil ekstraksi. Untuk
```

mengatasi masalah ini, Candès _et al_. (2011) mengembangkan metode statistik
baru bernama _Robust_ PCA (RPCA).
Eman _et al_. (2023) menjelaskan lebih lanjut bahwa RPCA
mendekomposisi data menjadi komponen-komponen utamanya, yang
merepresentasikan struktur dasar dari data, sekaligus mengidentifikasi dan
menghilangkan _outlier_ (pencilan) serta _noise_ (derau). Objek yang dijadikan
bahan penelitian dalam konteks ini adalah wajah manusia dan _outlier_ (pencilan)
yang berupa oklusi pada wajah seperti masker ataupun kacamata.
Berdasarkan penelitian oleh Eman _et al._ (2023), proses utama dari RPCA
adalah mendekomposisi sebuah matriks data masukan (𝑋) menjadi penjumlahan
dari dua matriks komponen yang terpisah, yaitu matriks _low-rank_ (𝐿) dan
matriks _sparse_ (𝑆).

```
𝑋=𝐿+𝑆
```
Dalam konteks pengenalan wajah, setiap matriks memiliki makna sebagai
berikut:

1. Matriks Masukan (X) merupakan matriks di mana setiap kolomnya
    adalah satu citra wajah yang telah diubah menjadi vektor piksel.
    Matriks ini merepresentasikan data asli yang mungkin mengandung
    gangguan.
2. Komponen _Low-Rank_ (L) merupakan representasi struktur data yang
    mendasar dan global, citra wajah yang bersih dan ideal.
3. Komponen _Sparse_ (S) merupakan representasi data _outlier_ (pencilan)
    atau _error_ besar yang lokasinya jarang ( _sparse_ ). Matriks ini
    menangkap oklusi seperti masker, kacamata, atau bayangan pekat yang
    ada pada citra wajah asli.
Secara matematis, tujuan dekomposisi ini diformulasikan sebagai sebuah
masalah optimisasi. Bentuk ideal dari masalah ini adalah mencari matriks L
dengan _rank_ serendah mungkin dan matriks S dengan elemen tak nol sesedikit
mungkin. Namun, masalah optimisasi ini bersifat _NP-hard_ , sehingga sangat sulit
untuk diselesaikan secara komputasi.
Upaya mengatasi permasalah optimisasi tersebut, digunakan pendekatan
relaksasi konveks ( _convex relaxation_ ) di mana fungsi yang sulit dioptimalkan
diganti dengan pendekatan terdekatnya yang bersifat konveks. Salah satu
algoritma yang paling terkenal untuk menyelesaikan masalah ini adalah
_Principal Component Pursuit_ (PCP) (Candès _et al_. 2011). Setelah proses
dekomposisi, matriks _low-rank_ L yang berisi citra wajah bersih kemudian dapat
digunakan untuk tahap ekstraksi fitur dan klasifikasi selanjutnya.

Optimisasi GA (Algoritma Genetika) berbasis SVM
Menurut Darmawan _et al_. (2023) Algoritma Genetika ( _Genetic Algorithm_ )
merupakan algoritma pencarian yang terinspirasi dari prinsip genetika dan
seleksi alam dari teori evolusi. Konsep utama dari algoritma ini adalah individu-
individu yang paling unggul akan bertahan hidup, sedangkan individu yang
lemah akan dieleminasi.
Dalam penerapan model pengenalan wajah, GA tidak digunakan untuk
mencari citra wajah terbaik, melainkan secara otomatis mencari kombinasi


```
parameter yang optimal untuk classifier SVM. Setiap “individu” dalam populasi
GA merepresentasikan satu set parameter SVM (seperti nilai C dan gamma).
"Individu" yang dianggap unggul adalah yang menghasilkan akurasi klasifikasi
tertinggi, sehingga proses evolusi ini pada akhirnya akan menghasilkan model
SVM dengan performa yang maksimal.
Support Vector Machine (SVM) merupakan suatu teknik machine learning
yang bertujuan untuk melakukan prediksi baik dalam kasus klasifikasi maupun
regresi. Teknik ini berusaha menemukan fungsi pemisah ( classifier ) terbaik
antara fungsi yang lain untuk nmemisahkan dua macam objek, pemisah tersebut
dinamankan hyperplane (Talib et al. 2024).
Berdasarkan penelitian oleh Putra et al. (2023), terdapat struktur secara
umum dari algoritma genetika, sebagai berikut:
```
1. Membangun populasi awal sekumpulan N individu secara acak.
    Individu dalam konteks ini adalah satu set kombinasi parameter SVM
    yang akan diuji (nilai C dan gamma)
2. Evaluasi kebugaran untuk menentukan tingkat kebugaran tiap individu
    dengan cara melatih model SVM, menguji performa SVM, dan
    menghasilkan nilai akurasi. Semakin tinggi akurasinya, semakin
    “bugar” individu tersebut.
3. Seleksi individu-individu dengan skor kebugaran tertinggi (yang
    menghasilkan akurasi SVM terbaik) dipilih sebagai "induk" untuk
    menciptakan generasi berikutnya. _Crossover_ untuk menghasilkan
    individu baru.
4. _Crossover_ melibatkan pertukaran "gen" (nilai parameter).
5. Mutasi untuk menjaga keragaman dan mencegah solusi prematur,
    beberapa individu anak akan mengalami perubahan kecil secara acak
    pada salah satu gen-nya (nilai parameternya)
Dengan struktur klasifikasi dan dekomposisi yang stabil dan teratur,
RPCA dan GA-SVM menjadi dua algoritma pilihan untuk berbagai kebutuhan
deteksi dan pengenalan objek, terutama citra wajah.

### 2.5 Arsitektur Raspberry Pi

Raspberry Pi merupakan sebuah _single-board computer_ tersusun yang dapat
mengoperasikan _operating system_ (OS) yang komprehensif (Aboluhom dan
Kandilli 2024). Dengan reputasi ketahanannya, perangkat ini berfungsi secara
terus-menurus, seperti mesin _server_ , dan mempertahankan konsumsi daya yang
minimal. Panas yang dihasilkan oleh CPU-nya dapat diabaikan, menjadikannya
pilihan efisien untuk berbagai aplikasi. Khususnya untuk pemakaian _machine
learning_ sistem presensi pengenalan wajah.

```
Penggunaan Raspberry Pi 4
Raspberry Pi memiliki Camera Serial Interface (CSI), sebuah port khusus
yang dirancang untuk menghubungkan modul kamera secara langsung. Ini
memungkinkan transfer data video beresolusi tinggi dengan latensi rendah, jauh
lebih efisien dibandingkan menggunakan kamera USB biasa. Untuk sistem
pengenalan wajah yang responsif, ini adalah fitur kunci. Sistem operasi
Raspberry Pi adalah sistem operasi berbasis Linux yang stabil dan didukung
```

penuh. Sistem operasi ini membantu kemudahan instalasi _library machine
learning_ , seperti OpenCV, Scikit-Learn, _requests_ , dan lain-lain.
Berdasarkan penelitian oleh Aboluhom dan Kandilli (2024),
pengintegrasian Raspberry Pi dan teknologi IoT akan memfasilitasi
pengumpulan dan pemrosesan data yang efisien, sehingga memungkinkan
pembaruan _real-time_ dan peningkatan skalabilitas. Hal ini menguatkan
pendukungan penggunaan Raspberry Pi untuk keperluan pengenalan wajah
sebagai sistem presensi. Berikut merupakan spesifikasi dan fitur dari Raspberry
Pi 4.

```
Tabel 1 Spesifikasi dan fitur Raspberry Pi 4
```
```
Spesifikasie Detail
Cores 4
Threads 4
Frekuensi prosesor 1.50 GHz
Penggunaan daya 2.7 W – 6.4 W
Jumlah memori 4 GB
Prosesor Rpi4 B
```
Dimensi (^85) cm x 49 cm x 1.8 cm
Sumber: Aboluhom dan Kandili (2024)
Gambar 1 _Pinout_ Raspberry Pi
Sumber: https://www.youngwonks.com/blog/Raspberry-Pi- 4 - Pinout
Arsitektur Pelatihan _Offline_ dan Inferensi _Online_
Implementasi _machine learning_ pada perangkat Raspberry Pi 4 mengikuti
arsitektur standar industri yang secara fundamental memisahkan dua fase utama,
yaitu pelatihan _offline_ dan inferensi _online_. Menurut Doyu _et al._ (2020), sebuah
model inferensi _machine learning_ umum yang telah dilatih secara konvensional
tidak dapat dijalankan dengan optimal pada perangkat IoT yang terbatas karena
sumber daya komputasi dari perangkat terbatas sehingga tidak memadai. Fase
pelatihan _offline_ adalah proses yang sangat intensif secara komputasi dan
memakan waktu banyak. Tahap ini akan menjalankan tugas-tugas berat seperti
ekstraksi fitur RPCA pada seluruh _dataset_ dan optimisasi GA-SVM yang iteratif
dijalankan pada _Virtual Private Server_ (VPS) melalui layanan _machine learning_.


```
Fase kedua adalah inferensi online , yaitu fase operasional dimana
Raspberry Pi tidak melakukan pelatihan ulang. Sebaliknya, Pi hanya bertugas
menggunakan artefak model yang sudah diproses secara offline untuk
menjalankan prediksi secara real-time. Proses inferensi ini, dalam konteks
sistem presensi yaitu memproyeksikan satu citra wajah ke basis RPCA dan
menjalankannya melalui model SVM, sangat ringan, cepat, dan efisien secara
komputasi. Arsitektur ini adalah kunci kelayakan proyek e mbedded machine
learning , karena memungkinkan model yang kompleks, yang dilatih di
lingkungan berdaya tinggi, untuk (diterapkan) pada perangkat berdaya rendah
yang hemat biaya (Warden dan Situnayake 2020).
```
### 2.6 Integrasi Sistem API Jira

Atlassian adalah sebuah perusahaan yang menciptakan dan menawarkan
beragam perangkat lunak dan produk yang dirancang bagi para _software developers_
dan tim untuk mengelola serta berkolaborasi dalam proyek mereka (Batskihh 2023).
Salah satu rangkaian produk yang banyak digunakan adalah _platform_ bernama Jira,
yang menyediakan kemampuan pelacakan isu ( _issue tracking_ ) dan manajemen
proyek. Isu Jira dapat memiliki berbagai bentuk tergantung pada bagaimana sebuah
tim menggunakan _platform_ tersebut. Namun, di Jira _Software_ , isu umumnya
berkaitan dengan item pekerjaan indivdual, seperti fitur utama, persyaratan
pengguna, dan _bug_ perangkat lunak.

```
Arsitektur API Jira
Arsitektur API Jira memungkinan sistem eksternal, seperti backend Golang
yang berjalan pada server, untuk berkomunikasi dan bertukar data dengan server Jira
secara terprogram. Interaksi ini didasarkan pada tiga komponen teknologi utama,
a. Autentikasi
Dalam penjaminan keamanan, setiap request ke API Jira harus
diautentikasi. Metode standar yang digunakan adalah autentikasi dasar
yang dikombinasikan dengan API Token. Token unik yang dibuat dari
akun Atlassian digunakan sebagai “kunci” untuk memvalidasi setiap
panggilan API.
b. Jira Query Language (JQL)
JQL adalah bahasa kueri khusus milik Jira, yang memungkinkan
pengguna untuk membangun kueri berdasarkan kriteria spesifik,
seperti tipe isu, status, assignee , dan custom fields (Garcia-Escribano
et al. 2025). Dalam penelitian ini, JQL digunakan untuk menyaring
task berdasarkan parameter seperti assignee (karyawan), status, dan
rentang waktu harian.
c. Pertukaran data dalam format JSON
Komunikasi antara backend Golang dan server Jira terjadi melalui
protokol HTTP. Backend berhasil mengirimkan kueri JQL sebagai
parameter dalam permintaan GET. Jika berhasil, server Jira akan
memberikan respon dalam bentuk format JSON.
```

## III METODE

### 3.1 Lokasi dan Waktu Proyek Akhir

Penelitian rancang bangun sistem presensi pengenalan wajah akan
diimplementasikan di Kantor PT. PNM Menara Caraka, Kuningan, Jakarta Selatan.
Penelitian dilakukan di dalam kantor karena kebutuhan alat digunakan secara
internal dan pengguna yang memakai adalah karyawan PT. PNM. Waktu penelitian
sistem presensi mulai di bulan November 2025 hingga Mei 2026. Waktu yang
direncanakan sudah termasuk pencetusan konsep hingga implementasi alat secara
sesungguhnya dan aktual di keberadaan dalam ( _interior_ ) kantor PT. PNM.

### 3.2 Teknik Pengumpulan Data dan Analisis Data

```
Teknik Pengumpulan Data
Teknik pengumpulan data yang digunakan pada projek akhir ini sesuai
dengan tujuan yang diharapkan oleh peneliti, yaitu sebagai berikut:
```
1. Studi literatur, studi ini bertujuan untuk mengumpulkan informasi
    terkait teknologi yang digunakan, berbagai macam model _machine_
    _learning_ , dokumentasi penggunaan dan perhitungan algoritma,
    perangkat keras yang sesuai (penggunaan kamera dan sistem _Internet_
    _of Things_ ) serta implementasi otomatisasi yang serupa.
2. Pengujian rancang bangun, menguji model mesin rancang bangun
    secara berkala dengan tujuan untuk mengevaluasi sistem _machine_
    _learning_ dan tampilan kamera untuk optimalisasi pengenalan wajah.
3. Observasi, memantau proses optimasi ekstraksi citra wajah _machine_
    _learning_ serta menghitung otomatisasi validasi pemasukan KPI
    realisasi _task_ karyawan PT. PNM untuk mendapatkan pemahaman
    faktor-faktor yang mempengaruhi efisiensi dan akurasi perhitungan.
4. Pengambilan dan analisis data, setelah pengujian rancang bangun dan
    observasi telah dilaksanakan, data yang dianalisis adalah tingkat
    akurasi dan _latency_ akhir citra wajah dan juga otomatisasi validasi
    pemasukan KPI realisasi _task_ karyawan PT. PNM.

```
Analisis Data
Analisis data pada penelitian proyek akhir ini dilakukan melalui pengujian
data secara analisis statistik deskriptif. Statistik deskriptif adalah statistik yang
digunakan untuk menganalisis data dengan cara mendeskripsikan atau
menggambarkan data yang telah terkumpul sebagaimana adanya tanpa
bermaksud membuat kesimpulan yang berlaku untuk umum atau generalisasi
(Sugiyono 2020). Statistik data utama yang diambil untuk penelitian proyek
akhir ini ada tiga, yaitu data akurasi sistem pengenalan wajah, waktu latency
pengenalan wajah secara real-time , dan validasi otomatisasi realisasi task pada
KPI karyawan PT. PNM.
Akurasi diukur dari tingkat kesuksesan sistem dalam mengenali wajah
yang terdaftar di dalam database. Data hasil pengujian sistem akan diolah ke
dalam bentuk Confusion Matrix. Menurut Fahmy (2022), Confusion Matrix
```

```
adalah tabel khusus yang digunakan dalam machine learning untuk
mendeskripsikan dan menilai kinerja sebuah model klasifikasi. Dari matriks
tersebut, akan dihitung metrik performa utama yaitu Akurasi, Presisi, dan Recall
untuk menilai keandalan kombinasi algoritma RPCA dan GA-SVM.
Latency diukur dari waktu yang diperlukan sistem untuk mendeteksi wajah
hingga keputusan identifikasi secara real-time. Parameter yang akan dihitung
adalah waktu rata-rata ( mean ), waktu minimum, waktu maksimum, dan standar
deviasi dari latensi untuk menentukan kelayakan implementasi real-time.
Setelah itu, otomatisasi pemasukan KPI realisasi task karyawan PT. PNM
perlu melakukan tahap validasi sebagai analisis data. Jumlah task aktual yang
ditugaskan oleh sebuah pengguna dan statusnya diselesaikan di Jira harus sama
dengan jumlah task yang sudah dicatat oleh database sistem presensi pengenalan
wajah. Analisis ini bertujuan untuk membuktikan bahwa alur otomatisasi
perhitungan realisasi task berjalan dengan aktual dan akurat.
```
### 3.3 Prosedur Kerja

Dalam pembuatan alat sistem pengenalan wajah berbasis algoritma RPCA
dan GA-SVM, terdapat prosedur yang wajib dilaksanakan untuk mencapai hasil
yang sesuai direncanakan. Prosedur tersebut digunakan sebagai pedoman
pengerjaan proyek akhir ini.

```
Gambar 2 Prosedur kerja
```
Prosedur penelitian ini diawali dengan studi literatur. Proses studi literatur
mengambil berbagai macam referensi, sumber materi dasar, hingga hasil observasi
penelitian sebelumnya yang terkait dengan topik dan judul penelitian. Proses ini
khusus untuk mematangkan referensi teoritis penelitian. Tahap selanjutnya adalah
mulai perancangan alat sistem pengenalan wajah.
Komponen utama yang digunakan pada tahap ini adalah Raspberry Pi 4,
kamera, dan layar LCD dengan fokus utama terhadap sistem _machine learning_ di
dalam sistem operasi Raspberry Pi. Selain merancang alat secara _hardware_ dan
_software_ , pengambilan data wajah diawal juga penting sebagai bahan pembelajaran
algoritma genetika berbasis SVM. Proses ini akan membutuhkan sekitar 30-40 citra
wajah per individu dengan berbagai ekspresi, sudut, dan sikap. Data wajah tersebut
kemudian akan melalui proses pelatihan model secara _offline_. Fitur dari _dataset_
wajah akan diekstraksi menggunakan RPCA, dan Algoritma Genetika (GA) akan
digunakan untuk mencari parameter optimal bagi _classifier_ SVM. Hasil dari tahap
ini adalah satu model SVM final yang sudah terlatih dan siap di- _deploy_.


Selanjutnya adalah tahap pengujian alat, dimana kamera yang tersambung ke
Raspberry Pi 4 akan melakukan inferensi secara _online_. Tampilan wajah tersebut
akan timbul di layar LCD dan algoritma RPCA akan ekstraksi citra wajah yang
diambil. Model SVM yang sudah dilatih akan menentukan identitas wajah paling
serupa dan mengeluarkan keputusan bahwa wajah tersebut adalah identitas sesuai
dengan wajah yang ingin presensi. Keputusan ini akan masuk ke dalam _database_
dan memanggil API Jira. Jika pengujian alat berhasil, tahap akan lanjut ke
implementasi secara sesungguhnya.
Tahap implementasi merupakan tahap penggunaan sistem presensi
pengenalan di kantor PT. Permodalan Nasional Madani. Tahap ini akan
menghasilkan data seperti tingkat akurasi dan kecepatan pengenalan wajah, dan
juga tingkat kinerja karyawan melalui panggilan API Jira. Ketiga data ini akan
dianalisis menggunakan teknik analisis statistik deskriptif, yang akan menghasilkan
pembahasan mengenai rumusan masalah yang sudah disusun.

### 3.4 Daftar Komponen

Komponen perangkat keras yang digunakan untuk penelitian berada pada
Tabel 2 beserta fungsi dan alasan penggunaan dalam sistem utama.

```
Tabel 2 Daftar komponen perangkat keras
```
```
Komponen Fungsi Alasan penggunaan
Adapter 5V/3A Sumber daya listrik
utama untuk komponen
perangkat keras, seperti
Raspberry Pi 4, kamera
web , dan LCD.
```
```
Raspberry Pi 4
membutuhkan voltase
5V dengan suplai
ampere sebesar 3A.
```
```
Kamera web Perangkat keras
penangkap gambar
wajah dengan tingkat
resolusi tinggi.
```
```
Komponen esensial
dalam pengambilan data
citra wajah individu
dengan kualitas tinggi.
Kabel mini-HDMI Kabel penghubung
antara Raspberry Pi 4 ke
kamera web sebagai
proyektor tampilan
antarmuka pengguna.
```
```
Proyeksi real-time dan
sebagai penghubung
sederhana antara kamera
dan layar LCD.
```
```
Micro SD 32 GB Tempat penyimpanan
sistem operasi, citra
wajah, dan perhitungan
analisis pengenalan
wajah.
```
```
Pemilihan jumlah
penyimpanan 32 GB
karena kebutuhan
komponen perangkat
lunak pengenalan wajah.
Raspberry Pi 4 Mini komputer utama
untuk menjalankan
semua proses machine
learning dan pengenalan
wajah.
```
```
Penggerak seluruh
komponen machine
learning dengan tingkat
mobilitas yang sangat
tinggi.
```

```
Casing Raspberry Pi 4 Pelindung dan tampilan
eksternal utama sistem
pengenalan wajah,
beserta tempat peletakan
komponen perangkat
keras lainnya.
```
```
Pelindung sederhana
untuk seluruh komponen
yang digunakan.
```
```
Kipas pendingin
Raspberry Pi 4
```
```
Komponen krusial dalam
pemeliharaan
lingkungan dan suhu di
sekitar Raspberry Pi 4.
Layar sentuh 7 inch Menampilkan secara
real-time yang ditangkap
oleh kamera web dan
status presensi.
```
```
Beraksi sebagai tampilan
antarmuka status
presensi.
```
Selain kebutuhan perangkat keras, sistem pengenalan wajah ini juga
membutuhkan komponen perangkat lunak. Tabel 3 menyertakan komponen-
komponen tersebut beserta fungsi dan alasan penggunaan dalam sistem.

```
Tabel 3 Daftar komponen perangkat lunak
```
```
Komponen Fungsi Alasan Penggunaan
Virtual Private
Server dari Google
Cloud Provider
```
```
Infrastruktur server virtual
utama yang beroperasi
penuh waktu untuk meng-
hosting seluruh ekosistem
sistem absensi, meliputi
basis data dan backend.
```
```
Memberikan hak akses
penuh serta sumber daya
komputasi terdedikasi yang
menjamin stabilitas
performa sistem saat
memproses beban kerja
berat.
OpenCV Library utama yang
digunakan untuk
pengambilan citra wajah
dan memprosesnya
sehingga dapat mendeteksi
wajah secara akurat.
```
```
OpenCV merupakan salah
satu library paling umum
untuk pengenalan wajah
yang menyediakan modul
deteksi wajah berbasis
DNN.
Scikit-learn Implementasi classifier
utama pada sistem, yaitu
Support Vector Machine
(SVM).
```
```
SVM merupakan classifier
tingkat atas yang memiliki
tingkat akurasi tinggi secara
efisien
NumPy Mengubah citra wajah
menjadi format matriks
numerik.
```
```
Efisien dalam operasi
matematis pada data
gambar.
RPCA Membersihkan data wajah
dari noise dan ekstraksi
fitur dari wajah yang
sudah dibersihkan menjadi
vektor numerik.
```
```
Membuat sistem lebih tahan
terhadap kondisi gambar
dunia nyata yang tidak
sempurna sehingga
meningkatkan akurasi
model.
```

PyGAD Mengimplementasikan
algoritma genetika (GA)
untuk mengoptimalkan
_classifier_ SVM.

Pustaka yang mudah
digunakan untuk
menggabungkan algortima
genetik (GA) dan _tuning_
atau mengoptimalkan SVM.
Jira Mengintegrasi sistem
presensi ke dalam
lingkungan Jira PT. PNM.

Lingkungan PT. PNM
sudah menggunakan
_platform_ Jira sebagai
_platform_ utama manajemen
pekerjaan,
PostgreSQL Sistem manajemen basis
data relasional (RDBMS)
untuk menyimpan dan
mengelola data terstruktur.

Teknologi ini
diimplementasikan karena
memiliki standar keandalan
tingkat tinggi (ACID
_compliant_ ).
MinIO S3 Storage Sebagai tempat
penyimpanan objek untuk
mengelola berkas tidak
terstruktur berskala besar,
khususnya data citra (foto)
wajah pegawai.

MinIO digunakan untuk
mencegah kelebihan beban
pada basis data utama,
sekaligus memberikan
performa unggah-unduh
gambar yang sangat cepat
dan efisien
Golang dengan
Fiber

```
Berkecepatan tinggi yang
memproses logika bisnis
dan menjembatani
komunikasi antara
hardware (Raspberry Pi),
basis data, dan dashboard
web
```
Memiliki kemampuan
_concurrency_ yang cukup
serta konsumsi memori
yang sangat rendah. Fiber
memastikan proses _routing_
API berjalan dengan latensi
yang minimal
Docker dan Docker
Compose

```
Platform kontainerisasi
yang digunakan untuk
mengemas layanan API
Golang dan layanan ML
beserta dependensinya.
```
Memastikan konsistensi
lingkungan eksekusi dari
pengembangan hingga
produksi (VPS), serta
mempermudah deployment.
React dengan Vite
dan Shadcn

```
React, Vite, dan Shadcn
digabungkan untuk
membangun antarmuka
pengguna interaktif
berbasis Single Page
Application (SPA) pada
antarmuka layar sentuh
karyawan dan dashboard
administrator sistem.
```
```
Penggunaan Vite, React,
dan Shadcn bertujuan
mempercepat kompilasi
kode, menciptakan navigasi
mulus, dan menyediakan
komponen UI profesional
yang mudah disesuaikan.
```

### 3.5 Desain Sistem

```
Skematik Rangkaian
```
```
Gambar 3 Skematik rangkaian
```
```
Rancangan skematik pada Gambar 3 mengilustrasikan rangkaian
perangkat keras yang digunakan untuk membangun sistem presensi
otomatis. Pada sistem tersebut, terdapat empat komponen hardware utama.
Pusat pemrosesan keseluruhan sistem berada di komponen Raspberry Pi 4.
Single-board computer ini krusial untuk proses utama inferensi dan
penghubungan sistem dengan internet, khususnya ke dalam database dan
platform Jira. Raspberry Pi 4 memiliki spesifikasi yang dapat memenuhi
kebutuhan menjalankan model menggunakan RPCA dan GA-SVM. Selain
memproseskan pembelajaran data, Raspberry Pi 4 juga dihubungkan dengan
dua komponen lainnya sebagai tampilan utama, yaitu kamera web dan layar
LCD touch screen.
Kamera web eksternal digunakan untuk menangkap gambar wajah
karyawan secara real-time. Kamera ini terhubung dengan port USB di
Raspberry Pi 4. Sebagai tampilan utama karyawan, sistem ini menggunakan
layar LCD touch screen. Layar ini berfungsi untuk menampilkan live feed
dari kamera, status presensi, dan feedback lainnya kepada karyawan.
Seluruh komponen yang terhubung ke Raspberry Pi 4 membutuhkan
sumber daya listrik utama. Sumber daya ini diperoleh dari adapter 5V/3A
yang akan menyuplai listrik ke sistem utama. Untuk mengumpuni
kebutuhan daya proses Raspberry Pi 4, dibutuhkan dua kipas pendingin
yang sambung langsung dengan Raspberry Pi 4. Seperti yang ditunjukkan
oleh kotak badan casing utama, komponen inti seperti Raspberry Pi 4 dan
```

```
layar LCD dirancang untuk ditempatkan bersama di dalam satu casing
untuk menciptakan perangkat yang ringkas dan terintegrasi.
```
_Flowchart_
Dalam perancangan sistem presensi ini, alur kerja operasional
digambarkan melalui dua diagram alir ( _flowchart_ ) esensial. Diagram alir ini
berfungsi untuk membedakan dua fase metodologi yang krusial dalam
pengembangan sistem, yaitu fase pelatihan data pada Gambar 4 dan fase
operasional _real-time_ pada Gambar 5.

```
Gambar 4 Diagram flowchart fase pelatihan data
```

Gambar 5 Diagram _flowchart_ fase _real-time_


Fase pertama, yang divisualisasikan dengan Gambar 4, merupakan proses
pelatihan data _offline_ yang bertujuan untuk membangun dan mengoptimalkan
pembelajaran mesin dari sistem. Proses ini diawali oleh pengumpulan data wajah,
dimana video singkat akan diambil selama 10 detik dari setiap karyawan. Video
yang diambil memiliki beberapa instruksi, seperti hadap atas, kanan, kiri, dan
bawah, agar terdapat variasi untuk model pembelajaran mesin. Dari video tersebut,
akan diekstraksi 30 hingga 50 gambar yang akan diproses. Gambar-gambar tersebut
akan diproses oleh _library_ OpenCV DNN untuk deteksi wajah dan tiap gambar akan
masuk ke tahap ekstraksi fitur dengan _Robust_ PCA. Algoritma ini akan mengubah
data menjadi numerik dan memisahkan antara fitur wajah yang esensial dengan
_noise_ visual, yang akan menghasilkan sebuah vektor fitur untuk setiap gambar.
Setelah setiap wajah memiliki _identifier_ unik masing-masing, proses akan
dilanjutkan ke optimasi GA-SVM. Algoritma genetika (GA) digunakan untuk
mencari pengaturan ( _hyperparameter_ ) terbaik bagi _classifier_ utama SVM. GA akan
menginisialisasi sebuah populasi yang terdiri dari berbagai kombinasi pengaturan
dan memulai _loop_ evolusi. Dalam _loop_ ini, setiap kombinasi pengaturan dievaluasi
kebugarannya ( _fitness_ ) berdasarkan seberapa akurat SVM yang dihasilkannya.
Melalui proses Seleksi, _Crossover_ , dan Mutasi yang berulang-ulang, GA secara
bertahap menemukan kombinasi pengaturan yang paling optimal. Setelah proses
_loop_ selesai, _hyperparameter_ terbaik akan diambil dan digunakan untuk pelatihan
mode final. Sebuah model SVM akhir dilatih menggunakan seluruh data vektor
fitur dan pengaturan terbaik tersebut. Terakhir, model RPCA dan model SVM yang
telah final dan teroptimasi ini disimpan ke dalam _file_ , siap untuk digunakan pada
fase kedua.
Fase kedua, yang direpresentasikan oleh Gambar 5, adalah proses _online_ yang
berjalan secara _real-time_ di Raspberry Pi untuk proses presensi harian. Proses ini
dimulai dengan inisialisasi sistem utama, dimana semua model terlatih (OpenCV
DNN, _Robust_ PCA, dan SVM) dimuat ke dalam memori, dan juga inisialisasi
_feedback_ audio dan koneksi ke API notifikasi Jira. Setelah sistem berhasil
inisialisasi, sistem akan masuk ke _loop_ pengenalan utama yang berjalan tanpa henti,
kecuali sistem dimatikan sumber daya listriknya. Di dalam _loop_ ini, detektor
OpenCV DNN akan terus memidai _frame_ video untuk mencari wajah. Jika wajah
terdeteksi, maka gambar tersebut akan diproses. Pertama, model _Robust_ PCA
mengekstraksi vektor fitur dari wajah tersebut, kemudian vektor tersebut akan
langsung diklasifikasi oleh SVM untuk menentukan identitasnya.
Jika wajah tidak dikenal, sistem akan memainkan audio _error_ dan _loop_ akan
kembali berulang. Jika wajah berhasil dikenal, proses presensi akan jalan. Proses
ini akan memerikan _timer cooldown_ untuk mencegah _spam_ presensi. Jika _cooldown_
sudah lolos, maka sistem akan memeriksa status internal pengguna, jika belum
tercatat hadir, sistem akan menjalankan _Clock In_ dan jika sudah tercatat hadir,
sistem akan menjalankan _Clock Out_. Kedua aksi ini akan mengirimkan notifikasi
ke Jira dan memainkan audio _feedback_ sesuai. Sistem sudah siap menjalankan
deteksi wajah berikutnya.


```
Rancangan Model 3D
Rancangan model tiga dimensi merupakan gambaran serupa realita bentuk
alat sistem presensi. Tampilan alatnya seperti di Gambar 6.
```
```
Gambar 6 Bentuk 3D alat presensi
```
```
Gambar 7 Peletakan kamera pada alat presensi
```
### 3.6 Matriks Rencana Proyek Akhir

```
Tabel 4 Jadwal pelaksanaan proyek akhir
```

# IV HASIL DAN PEMBAHASAN

**4.1 Lokasi dan Waktu Proyek Akhir**

Penelitian rancang bangun sistem presensi pengenalan wajah akan
diimplementasikan di Kantor PT. PNM Menara Caraka, Kuningan, Jakarta Selatan.
Penelitian dilakukan di dalam kantor karena kebutuhan alat digunakan secara
internal dan pengguna yang memakai adalah karyawan PT. PNM. Waktu penelitian
sistem presensi mulai di bulan November 2025 hingga Maret 2026. Waktu yang
direncanakan sudah termasuk pencetusan konsep hingga implementasi alat secara
sesungguhnya dan aktual di keberadaan dalam ( _interior_ ) kantor PT. PNM.

**4.2 Teknik Pengumpulan Data dan Analisis Data**

Teknik Pengumpulan Data
Teknik pengumpulan data yang digunakan pada projek akhir ini sesuai
dengan tujuan yang diharapkan oleh peneliti, yaitu sebagai berikut:
Studi literatur, studi ini bertujuan untuk mengumpulkan informasi terkait teknologi
yang digunakan, berbagai macam model _machine_


## DAFTAR PUSTAKA

Aboluhom AAA, Kandilli I. 2024. Face recognition using deep learning on
Raspberry Pi. _Computer Journal_. 67(10):3020–3030.
doi:10.1093/comjnl/bxae066.
Batskihh J. 2023. DevOps Approach in Software Development Using Atlassian Jira
Software. Tartu: Playtech Estonia OÜ.
Bora MA, Hanafie A, Haslindah A. 2025. _Ergonomi Kognitif: Optimalisasi
Interaksi Manusia Dan Sistem_. Sari M, Ahsani DM, editor. Padang: GET
PRESS INDONESIA. https://www.researchgate.net/publication/389102112.
Candès EJ, Li X, Ma Y, Wright J. 2011. Robust principal component analysis?
_Journal of the ACM_. 58(3):1–37. doi:10.1145/1970392.1970395.
Darmawan IPDW, Pradnyana GA, Pascima IBN. 2023. Optimasi Parameter
Support Vector Machine Dengan Algoritma Genetika Untuk Analisis
Sentimen Pada Media Sosial Instagram. _SINTECH (Science and Information
Technology) Journal_. 6(1):58-67. https://doi.org/10.31598.
Doyu H, Morabito R, Höller J. 2020. Bringing Machine Learning to the Deepest
IoT Edge with TinyML as-a-Service*. _IEEE IoT Newsletter_.
https://www.researchgate.net/publication/342916900.
Eman M, Mahmoud TM, Tarek AE-R. 2023. A Novel Hybrid Approach to Masked
Face Recognition using Robust PCA and GOA Optimizer. _Scientific Journal
for Damietta Faculty of Science_. 13(3):25–35.

# doi:10.21608/sjdfs.2023.222524.1117.

Fahmy MM. 2022. Confusion Matrix in Binary Classification Problems: A Step-
by-Step Tutorial. _Journal of Engineering Research_. 6(5).
https://doi.org/10.21608/erjeng.2022.274526.
Garcia-Escribano J, Carbajo A, Aranguren ME, Lopez-Novoa U. 2025. Using GPT
to build a Project Management assistant for Jira environments.
https://doi.org/10.48550/arXiv.2509.26014
Li L, Mu X, Li S, Peng H. 2020. A Review of Face Recognition Technology. _IEEE
Access_. 8:139110–139120. doi:10.1109/ACCESS.2020.3011028.
Meutia KI, Alqorrib Y, Fauzi A, Langi Y, Fauziah YN, Apriyanto W, Ramadhani
ZI. 2022. Pengaruh Usia Karyawan Dan Absensi Karyawan Terhadap Kinerja
Karyawan. _JEMSI: Jurnal Ekonomi Manajemen Sistem Informasi._ 3(6):674-
681. doi:10.31933/jemsi.v3i6.
Oloyede MO, Hancke GP, Myburgh HC. 2020. A review on face recognition
systems: recent approaches and challenges. _Multimed Tools Appl_. 79(37–
38):27891–27922. doi:10.1007/s11042- 020 - 09261 - 2.
Putra IBWBA, Astuti LG, Karyawati AE, Santiyasa W, Pramartha CRA, Astawa
GS. Diagnosis Penyakit Retinopati Diabetes Menggunakan SVM dengan
Optimasi Algoritma Genetika. _Jurnal Elektronik Ilmu Komputer Udayana_.
11(3):457-468. https://doi.org/10.24843/JLK.2023.v11.i03.p01
Rahma SA, Astuti SB. 2025. Pengaruh Ergonomi pada Kantor Teknik terhadap
Produktivitas Kerja Karyawan. _Lintas Ruang: Jurnal Pengetahuan dan
Perancangan Desain Interior_. 13(1):1–73.
doi:https://doi.org/10.24821/lintas.v13i1.15116.


Razaqa D, Al Maghribi MR, Gunasti N, Wati T. Analisis Etika dan Dampak
Penggunaan Sistem Pengenalan Wajah untuk Manajemen Kehadiran di
Lingkungan Sekolah. _Informatik Jurnal Ilmu Komputer_. 20(3):50-57.
https://doi.org/10.52958/iftk.v20i2.8093
Saied M, Syafii A. 2023. Perancangan dan Implementasi Sistem Absensi Berbasis
Teknologi Terkini Untuk Meningkatkan Efisiensi Pengelolaan Kehadiran
Karyawan dalam Perusahaan. _Jurnal Teknik Indonesia_. 2(3):87–92.
doi:10.58860/jti.v2i3.21.
Sari IP, Ramadhani F, Satria A, Apdilah D. 2023. Implementasi Pengolahan Citra
Digital dalam Pengenalan Wajah menggunakan Algoritma PCA dan Viola
Jones. _Hello World Jurnal Ilmu Komputer_. 2(3):146–157.
https://doi.org/10.56211/helloworld.v2i3.346.
Setianingrum HW. 2024. Urgensi Digitalisasi Manajemen Sumber Daya Manusia.
_Innovative: Journal of Social Science Research_. 4(3):10314– 10321.
https://doi.org/10.31004/innovative.v4i3.11668
Shakkak BIM, Al Mazruii SAKM. 2022. Face Recognition based on Convoluted
Neural Networks: Technical Review. _Applied computing Journal_. 2(2): 193 -

212. doi:10.52098/acj.202247.
Subiyanto S, Priliyana D, Riyadani ME, Iksan N, Wibawanto H. 2020. Face
recognition system with PCA-GA algorithm for smart home door security
using Rasberry Pi. _Jurnal Teknologi dan Sistem Komputer_. 8(3):210–216.
doi:10.14710/jtsiskom.2020.13590.
Sugiyono. 2020. _Metode Penelitian Kuantitatif, Kualitatif Dan R&D_. Ed ke-2.
Bandung: Penerbit Alfabeta.
Susilo AE, Abdurrahman. 2023. Manajemen Sumber Daya Manusia Untuk
Meningkatkan Kinerja Karyawan Melalui Absensi Digital. _Jurnal Educatio
FKIP UNMA_. 9(1):318–326. https://doi.org/10.31949/educatio.v9i1.4629.
Susim T, Darujati C, Artikel I. 2021. Pengolahan Citra Untuk Pengenalan Wajah
(Face Recognition) Menggunakan Opencv. _Jurnal Syntax Admiration_. 2(3):
534 - 545. https://doi.org/10.46799/jsa.v2i3.202
Talib S, Sudin S, Suratin MD. 2024. Penerapan Metode Support Vector Machine
(Svm) Pada Klasifikasi Jenis Cengkeh Berdasarkan Fitur Tekstur Daun. _Jurnal
Riset Sistem dan Teknologi Informasi (RESTIA)_. 2(1).
https://doi.org/10.30787/restia.v2i1.1364
Waliyaden FN, Leo GTW. 2024. Office Technology In The Digital Age: Reviewing
Modern Tools For Better Productivity And Choosing Ergonomic Office
Equipment For A Healthier Office. _PAJAMKEU : Pajak dan Manajemen
Keuangan_. 1(2):2–15.


