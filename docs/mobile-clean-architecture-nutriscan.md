# Mobile Clean Architecture Guide — NutriScan

## Context

NutriScan adalah aplikasi Flutter untuk preventive nutrition intelligence. MVP mobile berfokus pada tiga fitur utama:

1. AI Vision Scanner untuk scan foto makanan.
2. Preventive Nudge untuk memberi saran sederhana sebelum makan.
3. Simple Energy Trend untuk melihat tren estimasi asupan mingguan.

Struktur mobile harus dibuat agar mudah dibaca, mudah dikembangkan oleh tim, dan siap untuk integrasi backend serta AI inference service.

---

## Architecture Decision

Project Flutter menggunakan pendekatan:

> Feature-first Clean Architecture

Artinya, setiap fitur memiliki tiga layer utama sendiri:

```text
features/[feature_name]/
├─ data/
├─ domain/
└─ presentation/
```

Bukan menggunakan struktur global seperti:

```text
lib/
├─ data/
├─ domain/
└─ presentation/
```

Alasan memilih feature-first:

- Lebih mudah mencari file berdasarkan fitur.
- Cocok untuk MVP NutriScan yang fiturnya jelas: scan, nudge, trend, settings, user.
- Lebih mudah dikerjakan oleh beberapa developer.
- Lebih mudah dirawat saat fitur bertambah.
- Mengurangi kebiasaan lompat-lompat folder antar layer global.

---

## Target Folder Structure

Struktur akhir yang diinginkan:

```text
apps/mobile/lib/
├─ main.dart
│
├─ app/
│  ├─ app.dart
│  ├─ router/
│  │  └─ app_router.dart
│  ├─ theme/
│  │  └─ app_theme.dart
│  └─ di/
│     └─ injection.dart
│
├─ core/
│  ├─ constants/
│  │  └─ app_constants.dart
│  ├─ errors/
│  │  ├─ exceptions.dart
│  │  └─ failures.dart
│  ├─ network/
│  │  ├─ api_client.dart
│  │  └─ api_endpoints.dart
│  ├─ utils/
│  │  └─ result.dart
│  └─ extensions/
│     └─ .gitkeep
│
├─ shared/
│  ├─ widgets/
│  │  └─ .gitkeep
│  ├─ models/
│  │  └─ .gitkeep
│  └─ helpers/
│     └─ .gitkeep
│
├─ features/
│  ├─ scan/
│  │  ├─ data/
│  │  │  ├─ datasources/
│  │  │  │  └─ scan_remote_datasource.dart
│  │  │  ├─ models/
│  │  │  │  └─ scan_result_model.dart
│  │  │  └─ repositories/
│  │  │     └─ scan_repository_impl.dart
│  │  │
│  │  ├─ domain/
│  │  │  ├─ entities/
│  │  │  │  └─ scan_result.dart
│  │  │  ├─ repositories/
│  │  │  │  └─ scan_repository.dart
│  │  │  └─ usecases/
│  │  │     └─ analyze_food_image.dart
│  │  │
│  │  └─ presentation/
│  │     ├─ pages/
│  │     │  └─ scan_page.dart
│  │     ├─ controllers/
│  │     │  └─ scan_controller.dart
│  │     └─ widgets/
│  │        └─ .gitkeep
│  │
│  ├─ nudge/
│  │  ├─ data/
│  │  │  └─ .gitkeep
│  │  ├─ domain/
│  │  │  └─ .gitkeep
│  │  └─ presentation/
│  │     ├─ pages/
│  │     │  └─ .gitkeep
│  │     ├─ controllers/
│  │     │  └─ .gitkeep
│  │     └─ widgets/
│  │        └─ .gitkeep
│  │
│  ├─ trend/
│  │  ├─ data/
│  │  │  └─ .gitkeep
│  │  ├─ domain/
│  │  │  └─ .gitkeep
│  │  └─ presentation/
│  │     ├─ pages/
│  │     │  └─ .gitkeep
│  │     ├─ controllers/
│  │     │  └─ .gitkeep
│  │     └─ widgets/
│  │        └─ .gitkeep
│  │
│  ├─ settings/
│  │  ├─ data/
│  │  │  └─ .gitkeep
│  │  ├─ domain/
│  │  │  └─ .gitkeep
│  │  └─ presentation/
│  │     ├─ pages/
│  │     │  └─ .gitkeep
│  │     ├─ controllers/
│  │     │  └─ .gitkeep
│  │     └─ widgets/
│  │        └─ .gitkeep
│  │
│  └─ user/
│     ├─ data/
│     │  └─ .gitkeep
│     ├─ domain/
│     │  └─ .gitkeep
│     └─ presentation/
│        ├─ pages/
│        │  └─ .gitkeep
│        ├─ controllers/
│        │  └─ .gitkeep
│        └─ widgets/
│           └─ .gitkeep
│
└─ generated/
   └─ .gitkeep
```

---

## Layer Responsibility

### 1. Presentation Layer

Lokasi:

```text
features/[feature]/presentation/
```

Berisi:

- Page
- Widget
- Controller
- State management
- UI logic
- Form handling
- Loading, error, empty state

Contoh:

```text
features/scan/presentation/pages/scan_page.dart
features/scan/presentation/controllers/scan_controller.dart
features/scan/presentation/widgets/scan_result_card.dart
```

Aturan:

- Boleh memanggil use case dari domain.
- Tidak boleh langsung memanggil API client.
- Tidak boleh menyimpan logic bisnis utama.
- Tidak boleh parsing response API secara langsung.

---

### 2. Domain Layer

Lokasi:

```text
features/[feature]/domain/
```

Berisi:

- Entity
- Repository contract/interface
- Use case
- Business logic utama

Contoh:

```text
features/scan/domain/entities/scan_result.dart
features/scan/domain/repositories/scan_repository.dart
features/scan/domain/usecases/analyze_food_image.dart
```

Aturan:

- Harus pure Dart.
- Tidak boleh import Flutter Material/Cupertino.
- Tidak boleh bergantung pada data layer.
- Tidak boleh tahu detail API, JSON, HTTP, atau database.
- Menjadi pusat aturan bisnis fitur.

---

### 3. Data Layer

Lokasi:

```text
features/[feature]/data/
```

Berisi:

- Remote datasource
- Local datasource
- DTO/model response
- Repository implementation

Contoh:

```text
features/scan/data/datasources/scan_remote_datasource.dart
features/scan/data/models/scan_result_model.dart
features/scan/data/repositories/scan_repository_impl.dart
```

Aturan:

- Boleh memanggil API client.
- Boleh parsing JSON.
- Boleh mengubah model API menjadi entity domain.
- Mengimplementasikan repository contract dari domain.

---

## Global Folder Responsibility

### app/

Berisi konfigurasi aplikasi secara global:

```text
app/
├─ app.dart
├─ router/
├─ theme/
└─ di/
```

Digunakan untuk:

- Root widget aplikasi.
- Routing.
- Theme.
- Dependency injection.

---

### core/

Berisi hal teknis global yang dipakai banyak fitur:

```text
core/
├─ constants/
├─ errors/
├─ network/
├─ utils/
└─ extensions/
```

Digunakan untuk:

- API client global.
- Endpoint constants.
- Error/failure handling.
- Utility umum.
- Extension Dart umum.

---

### shared/

Berisi komponen yang reusable lintas fitur:

```text
shared/
├─ widgets/
├─ models/
└─ helpers/
```

Contoh:

- Custom button.
- Loading widget.
- Empty state widget.
- App card.
- Shared formatter/helper.

---

### generated/

Berisi file hasil code generation.

Aturan:

- Jangan isi manual kecuali memang hasil generator.
- Simpan `.gitkeep` agar folder tetap masuk Git.

---

## Import Direction Rules

Dependency harus mengarah ke dalam, bukan sebaliknya.

```text
presentation → domain
data → domain
data → core
presentation → shared
presentation → core
```

Yang tidak boleh:

```text
domain → presentation
domain → data
domain → Flutter UI
core → features
shared → features
```

---

## Feature Flow Example: Food Scan

Flow scan makanan:

```text
ScanPage
↓
ScanController
↓
AnalyzeFoodImage
↓
ScanRepository
↓
ScanRepositoryImpl
↓
ScanRemoteDataSource
↓
ApiClient
↓
Backend API / AI inference
```

Folder mapping:

```text
presentation/pages/scan_page.dart
presentation/controllers/scan_controller.dart
domain/usecases/analyze_food_image.dart
domain/repositories/scan_repository.dart
data/repositories/scan_repository_impl.dart
data/datasources/scan_remote_datasource.dart
core/network/api_client.dart
```

---

## Minimal Starter Files

### lib/main.dart

```dart
import 'package:flutter/material.dart';
import 'app/app.dart';

void main() {
  runApp(const NutriScanApp());
}
```

---

### lib/app/app.dart

```dart
import 'package:flutter/material.dart';
import 'theme/app_theme.dart';
import '../features/scan/presentation/pages/scan_page.dart';

class NutriScanApp extends StatelessWidget {
  const NutriScanApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'NutriScan',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      home: const ScanPage(),
    );
  }
}
```

---

### lib/app/theme/app_theme.dart

```dart
import 'package:flutter/material.dart';

class AppTheme {
  static ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      colorSchemeSeed: Colors.green,
      scaffoldBackgroundColor: const Color(0xFFF8FAF7),
    );
  }
}
```

---

### lib/features/scan/presentation/pages/scan_page.dart

```dart
import 'package:flutter/material.dart';

class ScanPage extends StatelessWidget {
  const ScanPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('NutriScan'),
      ),
      body: const Center(
        child: Text('AI Vision Scanner Page'),
      ),
    );
  }
}
```

---

## Refactor Plan

1. Buat branch baru:

```bash
git checkout -b chore/mobile-clean-architecture
```

2. Rapikan struktur folder `apps/mobile/lib`.

3. Pindahkan folder lama:

```text
lib/data/scan       → lib/features/scan/data
lib/domain/scan     → lib/features/scan/domain
lib/ui/scan         → lib/features/scan/presentation

lib/data/trend      → lib/features/trend/data
lib/domain/trend    → lib/features/trend/domain
lib/ui/trend        → lib/features/trend/presentation

lib/ui/nudge        → lib/features/nudge/presentation
lib/ui/settings     → lib/features/settings/presentation
lib/domain/user     → lib/features/user/domain
lib/data/api        → lib/core/network
```

4. Hapus folder layer global lama jika sudah kosong:

```text
lib/data
lib/domain
lib/ui
```

5. Pastikan `main.dart` memanggil `app/app.dart`.

6. Pastikan app tetap bisa jalan:

```bash
flutter analyze
flutter test
flutter run
```

7. Commit perubahan:

```bash
git add .
git commit -m "chore: restructure mobile app with feature-first clean architecture"
```

---

## Naming Convention

Gunakan snake_case untuk file dan folder.

Benar:

```text
scan_page.dart
scan_controller.dart
scan_repository.dart
scan_repository_impl.dart
analyze_food_image.dart
```

Hindari:

```text
ScanPage.dart
scanController.dart
RepositoryImpl.dart
```

Gunakan nama folder fitur yang singkat dan jelas:

```text
scan
nudge
trend
settings
user
```

---

## Notes for Future Development

- Mulai implementasi dari fitur `scan` karena ini fitur utama MVP.
- Jangan membuat abstraction berlebihan sebelum kebutuhan jelas.
- Untuk state management, siapkan folder `controllers` dulu. Library seperti Riverpod, Bloc, atau Provider bisa ditentukan kemudian.
- Untuk integrasi API, simpan konfigurasi global di `core/network`.
- Untuk AI inference, Flutter tidak menyimpan model utama. Flutter hanya mengirim gambar ke backend/API, lalu menerima hasil inference.
