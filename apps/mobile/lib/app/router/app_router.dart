import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../features/home/presentation/pages/home_page.dart';
import '../../features/scan/presentation/pages/scan_page.dart';
import '../../features/scan/presentation/pages/scan_result_page.dart';
import '../../features/trend/presentation/pages/trend_page.dart';
import '../../features/history/presentation/pages/history_page.dart';

final GlobalKey<NavigatorState> _rootNavigatorKey = GlobalKey<NavigatorState>();

class AppRouter {
  static final router = GoRouter(
    navigatorKey: _rootNavigatorKey,
    initialLocation: '/',
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const HomePage(),
      ),
      GoRoute(
        path: '/scan',
        builder: (context, state) => const ScanPage(),
      ),
      GoRoute(
        path: '/result',
        builder: (context, state) => const ScanResultPage(),
      ),
      GoRoute(
        path: '/trend',
        builder: (context, state) => const TrendPage(),
      ),
      GoRoute(
        path: '/history',
        builder: (context, state) => const HistoryPage(),
      ),
    ],
  );
}
