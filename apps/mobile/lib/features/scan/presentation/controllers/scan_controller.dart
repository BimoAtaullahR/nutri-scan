import 'package:flutter_riverpod/flutter_riverpod.dart';

enum ScanStatus { idle, analyzing, success, error }

class ScanState {
  final ScanStatus status;
  final String? errorMessage;

  ScanState({this.status = ScanStatus.idle, this.errorMessage});

  ScanState copyWith({ScanStatus? status, String? errorMessage}) {
    return ScanState(
      status: status ?? this.status,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }
}

class ScanController extends Notifier<ScanState> {
  @override
  ScanState build() {
    return ScanState();
  }

  void reset() {
    state = ScanState();
  }

  Future<void> analyzeImage() async {
    state = state.copyWith(status: ScanStatus.analyzing);
    await Future.delayed(const Duration(seconds: 2));
    state = state.copyWith(status: ScanStatus.success);
  }
}

final scanControllerProvider = NotifierProvider<ScanController, ScanState>(() {
  return ScanController();
});
