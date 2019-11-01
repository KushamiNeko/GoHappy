import "dart:html";

class IndexInput {
  final Element _container = querySelector("#index-input-container");

  final InputElement _idate = querySelector("#index-input-date");
  final InputElement _isymbol = querySelector("#index-input-symbol");
  final InputElement _ifreq = querySelector("#index-input-frequency");
  final InputElement _iversion = querySelector("#index-input-version");
  final ButtonElement _btn = querySelector("#index-input-button");

  bool _showRecords = false;

  IndexInput() {
    showRecords = false;
  }

  Element get container => _container;

  String get date => _idate.value;

  String get symbol => _isymbol.value;

  String get frequency => _ifreq.value;

  void setInputs(String value) {
    var regex = new RegExp(r"\s+");
    value = value.replaceAll(regex, " ");

    var vs = value.split(" ");

    if (vs.length < 3 || vs.length > 5) {
      throw new Exception("invalid value: ${value}");
    }

    regex = new RegExp(r"\d{8} \d{2}:\d{2}:\d{2}");
    if (regex.hasMatch(value)) {
      _idate.value = "${vs[0].trim()} ${vs[1].trim()}";
      _ifreq.value = vs[3].trim().toLowerCase();
      if (_showRecords) {
        _iversion.value = vs[4].trim();
      }
    } else {
      _idate.value = vs[0].trim();
      _ifreq.value = vs[2].trim().toLowerCase();
      if (_showRecords) {
        _iversion.value = vs[3].trim();
      }
    }

    // _isymbol.value = vs[1].trim().toLowerCase();
  }

  bool get showRecords => _showRecords;

  void set showRecords(bool ans) {
    _showRecords = ans;

    if (_showRecords) {
      _iversion.parent.classes.remove("index-input-text-hidden");
    } else {
      _iversion.parent.classes.add("index-input-text-hidden");
    }
  }

  String value() {
    String symbol;

    if (_isymbol.value.contains(",")) {
      symbol = _isymbol.value.split(",")[0];
    } else {
      symbol = _isymbol.value;
    }

    var val = "${_idate.value} ${symbol} ${_ifreq.value}";
    if (_showRecords) {
      val = "${val} ${_iversion.value}";
    }

    return val;
  }

  void activate(Function func(String value)) {
    _btn.onClick.listen((Event event) {
      func(value());
    });
  }

  void set date(String date) {
    _idate.value = date;
  }

  void set symbol(String symbol) {
    _isymbol.value = symbol;
  }

  void set frequency(String frequency) {
    _ifreq.value = frequency;
  }

  bool isFocus() {
    var focused = document.activeElement;
    return (focused == _idate ||
        focused == _isymbol ||
        focused == _ifreq ||
        focused == _iversion ||
        focused == _btn);
  }
}
