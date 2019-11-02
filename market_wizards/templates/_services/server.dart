import "dart:async";
import "dart:html";
import "dart:convert";

class Server {
  String _time;
  StreamController<String> _$time;

  String _frequency = "d";
  StreamController<String> _$frequency;

  String _symbol = "esz19";

  String _version = "1";
  StreamController<String> _$version;

  String _function = "refresh";

  bool _showRecords = false;
  StreamController<bool> _$showRecords;

  StreamController<String> _$chartUrl;

  //StreamController<String> _$chartInspect;
  StreamController<Map<String, dynamic>> _$chartInspect;

  static Server _server = null;

  bool _working = true;

  factory Server() {
    if (_server == null) {
      _server = Server._internal();
    }

    return _server;
  }

  Server._internal()
      : _$time = StreamController.broadcast(),
        _$frequency = StreamController.broadcast(),
        _$showRecords = StreamController.broadcast(),
        _$version = StreamController.broadcast(),
        _$chartUrl = StreamController.broadcast(),
        _$chartInspect = StreamController.broadcast() {
    var now = new DateTime.now();
    _time =
        "${now.year.toString()}${now.month.toString().padLeft(2, "0")}${now.day.toString().padLeft(2, "0")}";
    _$time.add(_time);
  }

  Stream get $time => _$time.stream;
  Stream get $frequency => _$frequency.stream;
  Stream get $version => _$version.stream;
  Stream get $showRecords => _$showRecords.stream;

  Stream get $chartUrl => _$chartUrl.stream;
  Stream get $chartInspect => _$chartInspect.stream;

  void broadcast() {
    _$time.add(_time);
    _$frequency.add(_frequency);
    _$version.add(_version);
    _$showRecords.add(_showRecords);
  }

  void recordsRequest(bool ans) {
    _function = "simple";

    _showRecords = ans;
    _$showRecords.add(_showRecords);

    getChart();
  }

  void forward() {
    _function = "forward";
    getChart();
  }

  void backward() {
    _function = "backward";
    getChart();
  }

  void symbolRequest(String symbol) {
    assert(new RegExp(r"^[a-zA-Z]{2,6}(?:\d{2})*$").hasMatch(symbol));
    _function = "simple";

    _symbol = symbol;

    getChart();
  }

  void freqRequest(String freq) {
    assert(new RegExp(r"h|d|w|m").hasMatch(freq));
    _function = "simple";

    _frequency = freq;
    _$frequency.add(freq);

    getChart();
  }

  void inputsRequest(String symbol, String time, String freq,
      {String version = "1"}) {
    assert(new RegExp(r"^[a-zA-Z]{2,6}(?:\d{2})*$").hasMatch(symbol));
    assert(new RegExp(r"h|d|w|m").hasMatch(freq));
    assert(new RegExp(r"^\d{8}$").hasMatch(time));

    _function = "refresh";

    _symbol = symbol;
    _frequency = freq;
    _time = time;

    if (_showRecords) {
      assert(new RegExp(r"^\d+$").hasMatch(version));
      _version = version;
    }

    getChart();
  }

  String _requestUrl() {
    var url =
        "${window.location.origin}/practice/${_symbol}/${_frequency}/${_function}/${_time}";

    if (_showRecords) {
      url = "${url}/records/${_version}";
    }

    url = "${url}?timestemp=${new DateTime.now().millisecondsSinceEpoch}";

    return url;
  }

  void timeRequest() async {
    if (_working) {
      return;
    }

    _function = "time";
    var url = _requestUrl();

    var info = await HttpRequest.getString(url);
    _time = info;

    _$time.add(info);
  }

  void inspectRequest(num x, num y) async {
    assert(x >= 0 && y >= 0);

    if (_working) {
      return;
    }

    _function = "inspect";
    var url = _requestUrl();

    url = "${url}&x=${x}&y=${y}";

    var info = await HttpRequest.getString(url);
    var m = json.decode(info);

    //_$chartInspect.add(info);
    _$chartInspect.add(m);
  }

  void getChart() {
    _working = true;

    var url = _requestUrl();

    _$chartUrl.add(url);
  }

  bool get isWorking => _working;

  void done() {
    _working = false;
  }
}
