# NutriScan UI Brief for Antigravity

## 0. Tujuan Brief

Brief ini digunakan untuk membantu Antigravity membangun UI Flutter NutriScan secara terarah, konsisten, dan maintainable. Fokus pekerjaan adalah **slicing UI dan pengondisian workspace frontend**, bukan implementasi backend/AI inference.

Aplikasi NutriScan adalah produk preventive nutrition intelligence. UI harus terasa ringan, sehat, bersih, dan cepat dipahami karena produk ini membantu pengguna mengambil keputusan sebelum makan, bukan membuat pengguna merasa sedang mengisi formulir tracking kalori yang rumit.

---

## 1. Konteks Produk

NutriScan berfokus pada kebiasaan pengguna yang sering memfoto makanan sebelum makan. Aplikasi perlu mengubah momen tersebut menjadi proses evaluasi cepat melalui:

- AI Vision Scanner untuk mengenali makanan dari foto.
- Visual Dominance Detection untuk membaca dominasi porsi secara visual.
- Preventive Nudge untuk memberi saran sederhana seperti “sisihkan porsi”.
- Simple Energy Trend untuk melihat perkembangan asupan energi mingguan.

Prinsip utama UI:

- Beban kognitif rendah.
- Feedback cepat dan langsung.
- Informasi utama terlihat tanpa banyak interaksi tambahan.
- Tidak terasa seperti aplikasi medis berat.
- Tidak menjanjikan diagnosis medis.

---

## 2. Catatan Penting dari PM / Desain

Pada desain homepage terdapat komentar berbentuk huruf **“G” berwarna pink**. Itu adalah **komentar dari PM**, bukan bagian dari UI.

Instruksi untuk Antigravity:

- Jangan buat widget, icon, badge, floating element, dekorasi, atau teks apa pun berdasarkan elemen “G” pink tersebut.
- Perlakukan elemen tersebut sebagai annotation/comment layer dari desain.
- Jika melakukan slicing dari screenshot/Figma, pastikan layer komentar/annotation tidak ikut masuk ke implementasi.

---

## 3. Analisis Desain dan Arah Visual

### 3.1 Personality UI

UI NutriScan sebaiknya memiliki karakter:

- Clean health-tech.
- Friendly dan approachable.
- Tidak terlalu klinis.
- Mobile-first.
- Fokus pada visual makanan, scanner, dan ringkasan insight.

Aplikasi ini bukan aplikasi diet ekstrem. Tone visualnya harus memberi kesan membantu pengguna mengambil keputusan yang lebih baik, bukan menghakimi makanan pengguna.

### 3.2 Hierarki Informasi

Prioritas informasi pada homepage:

1. Status/ringkasan harian pengguna.
2. CTA utama untuk scan makanan.
3. Ringkasan hasil atau rekomendasi preventif terakhir.
4. Trend energi mingguan.
5. Riwayat scan singkat atau insight tambahan.

Scanner harus menjadi aksi utama. Hindari homepage yang terlalu padat dengan banyak card setara, karena MVP membutuhkan alur yang cepat dari buka aplikasi → scan → lihat feedback.

### 3.3 Layout Direction

Gunakan layout berbasis card dengan spacing lapang:

- Background terang, clean, dan tidak ramai.
- Card utama memiliki radius besar.
- CTA scan dibuat paling menonjol.
- Komponen trend dibuat sederhana, bukan chart yang terlalu kompleks.
- Area nudge menggunakan bahasa singkat dan actionable.

Rekomendasi struktur visual homepage:

- Header greeting + subtitle ringan.
- Main scan card / camera CTA.
- Preventive insight card.
- Weekly energy trend card.
- Recent scans / meal history compact list.

---

## 4. Color Palette Recommendation

Catatan: palette ini adalah rekomendasi awal berdasarkan positioning NutriScan sebagai preventive nutrition app. Jika desain Figma sudah memiliki warna final, gunakan warna Figma sebagai source of truth dan mapping-kan ke token di bawah.

### 4.1 Core Palette

| Token | Rekomendasi Warna | Fungsi |
|---|---:|---|
| Primary Green | `#22C55E` | CTA utama, scan button, positive health action |
| Deep Green | `#166534` | Text emphasis, icon aktif, heading kecil |
| Mint Surface | `#ECFDF5` | Background card ringan, success/nutrition surface |
| Soft Lime | `#D9F99D` | Accent ringan untuk highlight sehat |
| Energy Orange | `#F97316` | Warning ringan, energy/calorie highlight, nudge |
| Soft Orange | `#FFEDD5` | Background nudge card |
| Neutral Dark | `#111827` | Heading utama |
| Neutral Body | `#4B5563` | Body text |
| Neutral Muted | `#9CA3AF` | Secondary text, placeholder |
| App Background | `#F8FAFC` | Page background |
| Card White | `#FFFFFF` | Surface utama |
| Border Soft | `#E5E7EB` | Divider dan card border |

### 4.2 Penggunaan Warna

- Gunakan hijau sebagai warna utama karena berasosiasi dengan health, freshness, dan prevention.
- Gunakan orange secara terbatas untuk menandai energi, kalori, atau nudge. Jangan terlalu banyak agar tidak terasa seperti warning berlebihan.
- Hindari merah sebagai warna dominan kecuali untuk error state, karena dapat membuat user merasa disalahkan.
- Gunakan background netral terang agar foto makanan dan CTA scanner tetap menjadi fokus.

### 4.3 Semantic Color Mapping

| Semantic Role | Warna |
|---|---:|
| Success / good balance | Primary Green |
| Preventive nudge | Energy Orange |
| Neutral information | Neutral Body |
| Disabled / inactive | Neutral Muted |
| Error | Red secukupnya, hanya untuk failure state |
| Background | App Background |
| Surface | Card White / Mint Surface |

---

## 5. Typography dan Spacing

### 5.1 Typography

Gunakan font default Flutter terlebih dahulu jika belum ada font khusus dari desain. Hindari terlalu banyak variasi ukuran.

Rekomendasi hierarchy:

- Display / Hero title: untuk greeting atau headline scan.
- Section title: untuk judul card seperti “Weekly Energy Trend”.
- Body: untuk penjelasan singkat.
- Caption: untuk metadata seperti waktu scan atau status kecil.
- Button label: medium/semi-bold.

### 5.2 Spacing

Gunakan spacing konsisten:

- Page horizontal padding: 20–24.
- Card padding: 16–20.
- Gap antar section: 20–28.
- Gap antar item kecil: 8–12.
- Radius card: 20–28.
- Radius button utama: 16–24.

Jangan membuat UI terlalu rapat karena produk ini harus terasa ringan dan low cognitive load.

---

## 6. Komponen UI yang Perlu Dibangun

Implementasi dilakukan **per komponen**, bukan langsung membangun screen besar sekaligus. Setiap komponen harus reusable dan tidak menyimpan logic berat di dalam widget.

### 6.1 Foundation Components

- AppScaffold
- AppHeader
- AppCard
- AppButton
- AppIconButton
- AppTextField jika dibutuhkan
- AppBottomNavigation
- AppLoadingState
- AppEmptyState
- AppErrorState

### 6.2 Feature Components: Home

- HomeGreetingHeader
- DailySummaryCard
- ScanFoodCtaCard
- PreventiveNudgeCard
- WeeklyEnergyTrendCard
- RecentScanList
- RecentScanItem

### 6.3 Feature Components: Scanner

- ScannerPreviewFrame
- ScanInstructionCard
- CaptureButton
- ImagePreviewCard
- ScannerLoadingOverlay
- ScanResultSummaryCard
- VisualDominanceCard
- PortionSuggestionCard

### 6.4 Feature Components: Trend / Dashboard

- EnergyTrendCard
- WeeklyBarSummary atau simple chart component
- EnergyInsightCard
- TrendFilterChip

### 6.5 Feature Components: History

- MealHistoryList
- MealHistoryItem
- MealTagChip
- MealDetailHeader
- MealNutritionEstimateCard

---

## 7. Screen Scope untuk MVP UI

Prioritaskan screen berikut:

### 7.1 Home Screen

Tujuan:

- Memberikan entry point cepat untuk scan.
- Menampilkan ringkasan energi/trend sederhana.
- Menampilkan preventive insight terakhir.

Komponen utama:

- HomeGreetingHeader
- ScanFoodCtaCard
- PreventiveNudgeCard
- WeeklyEnergyTrendCard
- RecentScanList

Catatan:

- Elemen “G” pink pada desain homepage harus diabaikan karena merupakan komentar PM.

### 7.2 Scanner Screen

Tujuan:

- Memberi pengalaman scan makanan yang cepat dan sederhana.
- Menampilkan state loading dan preview secara jelas.

Komponen utama:

- ScannerPreviewFrame
- ScanInstructionCard
- CaptureButton
- ScannerLoadingOverlay
- ImagePreviewCard

Catatan:

- Untuk tahap UI, boleh gunakan mock image dan mock result.
- Jangan implementasikan AI inference asli dulu.

### 7.3 Result Screen

Tujuan:

- Menampilkan hasil estimasi makanan secara mudah dipahami.
- Memberikan nudge preventif yang actionable.

Komponen utama:

- ScanResultSummaryCard
- VisualDominanceCard
- PortionSuggestionCard
- EnergyEstimateCard
- SaveResultButton

Catatan:

- Hindari klaim medis.
- Pakai copy yang bersifat estimasi, misalnya “perkiraan”, “indikasi”, atau “saran ringan”.

### 7.4 Dashboard / Trend Screen

Tujuan:

- Menampilkan Simple Energy Trend mingguan.
- Membantu pengguna melihat pola tanpa input manual rumit.

Komponen utama:

- WeeklyEnergyTrendCard
- EnergyInsightCard
- RecentScanList

### 7.5 History Screen

Tujuan:

- Menampilkan riwayat scan makanan.
- Memberi akses ke detail hasil scan sebelumnya.

Komponen utama:

- MealHistoryList
- MealHistoryItem
- MealTagChip

---

## 8. Workspace Conditioning: Riverpod dan GoRouter

Project Flutter sudah dibuat dan struktur folder clean architecture sudah disiapkan berdasarkan guide internal. Antigravity harus mengikuti struktur yang sudah ada, bukan membuat struktur baru dari nol.

### 8.1 Prinsip Umum

- Jangan mengubah struktur folder besar tanpa alasan kuat.
- Ikuti guide internal yang sudah ada di repository.
- Gunakan Riverpod untuk state management.
- Gunakan GoRouter untuk routing.
- Pisahkan UI, state, routing, dan model mock.
- Jangan menaruh logic navigation kompleks langsung di widget kecil.
- Jangan melakukan integrasi backend/AI dulu kecuali sudah ada instruction lanjutan.

### 8.2 Dependency Conditioning

Pastikan dependency berikut tersedia di project:

- flutter_riverpod
- go_router
- optional: riverpod_annotation dan build_runner jika project memang sudah memakai code generation

Jika dependency belum ada, tambahkan secara minimal dan jangan mengubah dependency yang tidak relevan.

### 8.3 Riverpod Usage Rules

Gunakan Riverpod untuk:

- UI state scanner.
- Selected tab bottom navigation jika diperlukan.
- Mock scan result state.
- Loading/error/success state untuk simulasi scan.
- Filter atau selected period pada dashboard.

Hindari:

- Business logic berat di widget.
- Global mutable state biasa.
- Provider yang terlalu besar untuk semua fitur.
- Menyatukan state home, scanner, result, dan dashboard dalam satu provider besar.

Rekomendasi pendekatan:

- Buat provider per feature.
- State sederhana untuk MVP UI.
- Gunakan mock data lokal untuk preview UI.
- Pastikan state mudah diganti ke data asli saat integrasi backend/AI nanti.

### 8.4 GoRouter Usage Rules

Gunakan GoRouter untuk route utama:

- Home
- Scanner
- Scan Result
- Dashboard / Trend
- History

Routing harus:

- Punya nama route yang konsisten.
- Tidak hardcode string route berulang-ulang di banyak widget.
- Mendukung navigasi dari Home ke Scanner.
- Mendukung navigasi dari Scanner ke Result menggunakan mock result atau parameter sederhana.
- Mendukung bottom navigation jika desain mengarah ke tab-based layout.

Untuk MVP, route guard/auth belum perlu dibuat kecuali sudah ada requirement baru.

---

## 9. Rekomendasi Struktur Folder Feature-Based

Ikuti struktur yang sudah ada. Jika perlu menambahkan folder, gunakan pola feature-based agar mudah maintain.

Rekomendasi umum:

- core
  - theme
  - router
  - constants
  - widgets
  - utils
- features
  - home
    - presentation
    - providers
    - widgets
  - scanner
    - presentation
    - providers
    - widgets
  - result
    - presentation
    - providers
    - widgets
  - dashboard
    - presentation
    - providers
    - widgets
  - history
    - presentation
    - providers
    - widgets

Catatan:

- Jangan membuat folder data/domain terlalu kompleks jika untuk tahap ini hanya slicing UI.
- Jika clean architecture di guide internal sudah mewajibkan domain/data/presentation, ikuti guide tersebut.
- Untuk UI mock, cukup letakkan model mock secara rapi agar nanti mudah diganti dengan repository asli.

---

## 10. UI State yang Harus Disiapkan

Minimal state untuk MVP UI:

### 10.1 Scanner State

- idle
- cameraReady / ready
- imageCaptured
- analyzing
- success
- error

### 10.2 Result State

- foodName / detectedFoodLabel
- confidenceEstimate
- estimatedEnergy
- dominantPortionLabel
- suggestedAction
- saved / unsaved

### 10.3 Dashboard State

- selectedWeek
- weeklyEnergyData
- insightSummary
- loading / empty / error

Gunakan mock data agar UI dapat diuji tanpa backend.

---

## 11. Copywriting Direction

Gunakan bahasa Indonesia yang ringan, jelas, dan tidak menghakimi.

Contoh tone:

- “Scan makananmu sebelum makan.”
- “Perkiraan energi hari ini.”
- “Porsi terlihat cukup dominan.”
- “Coba sisihkan sedikit porsi karbo agar lebih seimbang.”
- “Trend minggu ini masih stabil.”

Hindari:

- “Makanan ini buruk.”
- “Kamu kelebihan kalori.”
- “Diagnosis risiko diabetes.”
- “Wajib kurangi makan.”

Gunakan kata seperti:

- perkiraan
- saran
- indikasi
- bantu
- seimbang
- ringan

---

## 12. Implementation Boundaries

Antigravity hanya perlu mengerjakan:

- Setup/pengkondisian Riverpod dan GoRouter.
- Slicing UI berdasarkan desain.
- Pembuatan komponen reusable.
- Mock data untuk kebutuhan preview.
- State UI sederhana untuk loading/success/error.
- Penyesuaian theme, color token, spacing, dan typography.

Antigravity tidak perlu mengerjakan:

- Integrasi AI model asli.
- API backend asli.
- Database production.
- Authentication.
- Medical recommendation engine.
- Perhitungan nutrisi klinis presisi.
- Upload image production-ready.

---

## 13. Execution Order untuk Antigravity

Kerjakan secara bertahap:

1. Baca guide struktur folder internal yang sudah ada.
2. Cek dependency Flutter saat ini.
3. Tambahkan atau rapikan Riverpod dan GoRouter secara minimal.
4. Buat theme token: color, typography, spacing, radius.
5. Buat foundation components.
6. Slicing Home Screen.
7. Slicing Scanner Screen dengan mock state.
8. Slicing Result Screen dengan mock data.
9. Slicing Dashboard/Trend Screen.
10. Slicing History Screen jika scope waktu cukup.
11. Refactor widget besar menjadi komponen kecil.
12. Pastikan tidak ada layer komentar desain seperti “G” pink yang masuk ke UI.
13. Jalankan format dan static analysis.

---

## 14. Acceptance Criteria UI

UI dianggap selesai untuk tahap ini jika:

- Home screen mengikuti desain dan tidak memasukkan komentar “G” pink sebagai UI.
- Color palette konsisten dan sudah dipetakan ke theme/token.
- Komponen dibuat modular, bukan satu file screen yang terlalu besar.
- Riverpod sudah digunakan untuk state UI yang relevan.
- GoRouter sudah digunakan untuk navigasi antar screen utama.
- Scanner memiliki mock flow: idle → analyzing → result/error.
- Result screen menampilkan estimasi dan preventive nudge secara jelas.
- Dashboard menampilkan simple weekly energy trend.
- Tidak ada logic AI/backend production di UI layer.
- Tidak ada hardcoded style berulang yang seharusnya masuk theme.
- Flutter analyze tidak menghasilkan error kritis.

---

## 15. Prompt Siap Pakai untuk Antigravity

Gunakan prompt berikut ke Antigravity:

> Bantu bangun UI Flutter untuk project NutriScan berdasarkan desain yang ada. Jangan langsung membuat implementasi backend/AI. Fokus pada slicing UI, reusable components, theme token, Riverpod untuk UI state, dan GoRouter untuk routing. Project Flutter sudah disetup dan struktur folder clean architecture sudah ada, jadi ikuti guide internal repository dan jangan membuat struktur besar baru tanpa alasan. Pada homepage ada komentar huruf “G” berwarna pink dari PM; jangan jadikan itu bagian dari UI. Bangun secara bertahap per komponen: foundation components, home, scanner, result, dashboard/trend, dan history. Gunakan mock data untuk scanner result dan trend. Pastikan UI clean, health-tech, low cognitive load, dan konsisten dengan color palette hijau/mint/orange sebagai aksen. Gunakan Context7 untuk mengecek dokumentasi terbaru Flutter, Riverpod, dan GoRouter bila diperlukan. Setelah selesai, jalankan format dan analyze, lalu laporkan file yang dibuat/diubah serta alasan perubahan.

---

## 16. Notes for Future Integration

Agar integrasi AI/BE nanti lebih mudah:

- Mock scan result harus menyerupai struktur data yang nanti dibutuhkan backend.
- State scanner jangan terlalu terikat pada UI tertentu.
- Route result harus bisa menerima data scan dari scanner.
- Komponen result harus bisa menerima data dari provider/repository di masa depan.
- Copy harus tetap menyebut estimasi, bukan angka klinis mutlak.

