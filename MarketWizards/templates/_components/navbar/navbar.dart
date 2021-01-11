import "dart:html";
import "../../_services/server.dart";

class Navbar {
  final DivElement _container;

  final ButtonElement _practice;
  final ButtonElement _records;

  final Server _server;

  final String _cls;

  bool _showRecords = false;

  Navbar(String id, {String cls = ""})
      : _cls = cls,
        _server = new Server(),
        _container = querySelector("#${id}-navbar-container"),
        _practice = querySelector("#${id}-navbar-practice"),
        _records = querySelector("#${id}-navbar-records") {
    var path = window.location.pathname;

    if (path.contains("practice")) {
      _practice.classes.add("${cls}-navbar-button-active");
    }

    _practice.onClick.listen((Event event) {
      window.location.pathname = "/view/practice";
    });

    _records.onClick.listen((Event event) {
      _showRecords = !_showRecords;

      if (_showRecords) {
        _records.classes.add("${cls}-navbar-button-active");
      } else {
        _records.classes.remove("${cls}-navbar-button-active");
      }

      _server.recordsRequest(_showRecords);
      _records.blur();
    });
  }

  void enterFullScreen(bool ans) {
    if (ans) {
      _container.classes.add("${_cls}-navbar-hidden");
    } else {
      _container.classes.remove("${_cls}-navbar-hidden");
    }
  }
}
