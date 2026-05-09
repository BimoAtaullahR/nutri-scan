
## Product Requirements Document (PRD)

# NutriScanAI

## 1. Executive Summary

```
Project ini bertujuan untuk membangun NutriScan AI, sebuah sistem preventive nutrition
intelligence yang dirancang untuk mengubah tren sosial "camera eats first" menjadi sistem
pendukung keputusan kesehatan proaktif. Fokus utamanya adalah memitigasi risiko penyakit
tidak menular, khususnya obesitas dan diabetes, dengan mengorkestrasi intervensi asupan
energi tepat di titik konsumsi sebelum pengguna mulai makan. Project ini dirancang dengan
prototipe 3 bulan untuk menghasilkan MVP yang fungsional
```
## 2. Product Description

```
Masyarakat Asia Pasifik saat ini menghadapi lonjakan penyakit tidak menular yang dipicu
oleh konsumsi energi berlebih tanpa disadari. Keputusan makan sehari-hari sering terjadi secara
impulsif tanpa sistem yang membantu mengevaluasi porsi. Aplikasi pencatat kalori konvensional
sangat merepotkan, dan kesadaran biasanya baru muncul setelah selesai makan. NutriScan AI
hadir untuk memberikan intervensi tepat waktu melalui Visual Inference dan Preventive Nudges ,
membantu pengguna mengevaluasi komposisi nutrisinya hanya dengan satu jepretan foto.
```
## 3. Product Goals

- Memperkenalkan sistem pendukung keputusan kesehatan preventif kepada masyarakat.
- Menurunkan risiko penyakit tidak menular (diabetes & obesitas) sesuai dengan
    penyelarasan SDG 3 (Target 3.4).
- Menyediakan intervensi berbasis estimasi energi instan yang memiliki beban kognitif
    rendah.
- Menyediakan visualisasi _Simple Energy Trend_ untuk melacak perkembangan asupan
    secara mingguan.

## 4. Success Metrics

- **Pre-Portioning Compliance Rate:** Target >30%, mengukur sejauh mana pengguna
    mematuhi saran AI untuk menyisihkan makanan.
- **System Latency:** Target < 3 detik dari proses jepretan hingga memunculkan _feedback_
    agar tidak mengganggu ritme makan.
- **30-Day Retention Rate:** Target >25% untuk mengukur keberlanjutan penggunaan dalam
    rutinitas harian.
- **Estimated Energy Intake Prevented:** Estimasi kalori berlebih yang berhasil dihindari
    melalui intervensi preventif.


## 5. Product Scope

```
In-Scope:
● Halaman AI Vision Scanner untuk pengenalan 50 kategori makanan lokal populer.
● Visual Dominance Detection untuk analisis komposisi porsi visual.
● Preventive Nudge Engine untuk notifikasi instan "Sisihkan Porsi".
● Simple Energy Trend dashboard
Out-Scope:
● Diagnosis atau rekomendasi medis klinis.
● Input perhitungan makronutrisi manual dalam hitungan gram pasti.
●
```
## 6. Key Requirements

- Fitur MVP wajib mencakup kemampuan deteksi dominasi visual dan mesin _nudge_
    preventif.
- Alur penggunaan aplikasi harus sangat intuitif dengan hierarki informasi, tipografi, dan
    _spacing_ yang konsisten agar beban kognitif pengguna tetap rendah.
- Aplikasi harus responsif dan dibangun dengan arsitektur _frontend_ modern (seperti Next.js
    dan Tailwind CSS) untuk memastikan _system latency_ tetap di bawah 3 detik.
- Dukungan infrastruktur _database_ dan pemrosesan _AI Inference_ yang dapat diukur
    efisiensinya.
-

## 7. Demographics

- **Target Pengguna:** Pengguna aktif dengan kebiasaan memfoto makanan ("camera eats
    first") serta pengguna B2B ( _corporate wellness_ & asuransi)
- **Lingkup:** Asia Pasifik (fokus awal pengenalan makanan lokal).
- **Edukasi:** Pengguna _smartphone_ yang membutuhkan laporan kesehatan tanpa proses
    _tracking_ yang merepotkan.
-
-

## 8. Constraints & Risks

```
Constraints :
● Waktu penyelesaian sangat singkat (30 hari).
● Keterbatasan pengalaman tim pengembang dengan beberapa tech stack spesifik.
Risks :
● Potensi keterlambatan dalam penyediaan data konten atau aset gambar produk
● Potensi latensi pemrosesan gambar jika server AI mengalami beban tinggi ( bottleneck ).
●
```
## 9. Timeline (Agile 4-Week Sprints)

### -


```
● Sprint 1 (Discovery & Blueprint): Fokus pada riset perilaku pengguna ("camera eats
first"), pembuatan alur kerja ( user-flow ), dan desain high-fidelity untuk fitur inti (Kamera &
Dashboard).
● Sprint 2 (Core Development): Fokus pada slicing UI dan implementasi fitur prioritas
tertinggi seperti fungsionalitas pemindaian gambar dan tracking mingguan.
● Sprint 3 (Intelligence & Integration): Fokus pada integrasi model AI Vision Scanner
(pengenalan 50 kategori makanan), Visual Dominance Detection , dan logika Preventive
Nudge Engine.
● Sprint 4 (Refinement & Launch): Fokus pada pengujian performa ( system latency <
detik), perbaikan bug , serta peluncuran aplikasi MVP ( Deployment ).
●
```
## 10. Product Deliverables

- **NutriScan Mobile/Web Build:** Akses ke _testing environment_ MVP.
- **Design System:** Dokumentasi komponen UI di Figma, termasuk _auto layout_ yang sudah
    melalui tahapan _review_ visual.
- **Technical Documentation:** Repositori kode untuk antarmuka serta arsitektur AI
    _inference_.
-

## 11. Project Team

```
● Gilbard (Hustler): Mengelola prioritas fitur, kemitraan B2B, dan strategi pitching.
● Hacker FE: Eksekusi teknis antarmuka aplikasi menggunakan framework modern yang
responsif.
● Hipster: Merancang antarmuka visual yang user-friendly dan memastikan konsistensi
desain untuk menekan beban kognitif pengguna.
● Hacker BE/AI: Mengelola database dan arsitektur AI inference agar respon gambar
stabil di bawah 3 deti
```
-


