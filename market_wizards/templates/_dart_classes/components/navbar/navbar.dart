import "dart:html";

class Navbar {
  final ButtonElement _study = querySelector("#navbar-study");
  final ButtonElement _practice = querySelector("#navbar-practice");
  final ButtonElement _records = querySelector("#navbar-records");

  bool _showRecords = false;

  Navbar() {
    var path = window.location.pathname;

    if (path.contains("study")) {
      _study.classes.add("navbar-button-active");
    }

    if (path.contains("practice")) {
      _practice.classes.add("navbar-button-active");
    }

    //_study.onClick.listen((Event event) {
    //window.location.pathname = "/view/study";
    //});

    _practice.onClick.listen((Event event) {
      window.location.pathname = "/view/practice";
    });

    //_records.onClick.listen((Event event) {
    //if (_showRecords) {
    //_records.classes.add("navbar-button-active");
    //} else {
    //_records.classes.remove("navbar-button-active");
    //}
    //});
  }

  void activate(void func(bool ans)) {
    _records.onClick.listen((Event event) {
      _showRecords = !_showRecords;

      if (_showRecords) {
        _records.classes.add("navbar-button-active");
      } else {
        _records.classes.remove("navbar-button-active");
      }

      func(_showRecords);
    });
  }
}
