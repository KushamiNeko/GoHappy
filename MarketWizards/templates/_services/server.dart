import "dart:async";
import "dart:html";
import "dart:convert";

class Server {
  String _time;
  StreamController<String> _$time;

  String _frequency = "d";
  StreamController<String> _$frequency;

  String _symbol = "es";
  StreamController<String> _$symbol;

  String _book = "es_${new DateTime.now().year - 1}";
  StreamController<String> _$book;

  String _note = "";
  StreamController<String> _$note;

  String _function = "refresh";

  bool _showRecords = false;
  StreamController<bool> _$showRecords;

  StreamController<String> _$chartUrl;

  StreamController<String> _$chartInspect;

  StreamController<Map<String, dynamic>> _$info;

  static Server _server = null;

  bool _working = false;

  factory Server() {
    if (_server == null) {
      _server = Server._internal();
    }

    return _server;
  }

  Server._internal()
      : _$time = StreamController.broadcast(),
        _$symbol = StreamController.broadcast(),
        _$frequency = StreamController.broadcast(),
        _$showRecords = StreamController.broadcast(),
        _$book = StreamController.broadcast(),
        _$note = StreamController.broadcast(),
        _$chartUrl = StreamController.broadcast(),
        _$chartInspect = StreamController.broadcast(),
        _$info = StreamController.broadcast() {
    var now = new DateTime.now();
    _time =
        "${now.year.toString()}${now.month.toString().padLeft(2, "0")}${now.day.toString().padLeft(2, "0")}";
    _$time.add(_time);
  }

  Stream get $time => _$time.stream;
  Stream get $symbol => _$symbol.stream;
  Stream get $frequency => _$frequency.stream;
  Stream get $book => _$book.stream;
  Stream get $note => _$note.stream;
  Stream get $showRecords => _$showRecords.stream;

  Stream get $chartUrl => _$chartUrl.stream;
  Stream get $chartInspect => _$chartInspect.stream;
  Stream get $info => _$info.stream;

  void broadcast() {
    _$time.add(_time);
    _$symbol.add(_symbol);
    _$frequency.add(_frequency);
    _$book.add(_book);
    _$note.add(_note);
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
    //assert(new RegExp(r"^[a-zA-Z]{2,6}(?:\d{2})*$").hasMatch(symbol));
    assert(new RegExp(r"^[a-zA-Z]{1,6}(?:\d{1,2})*$").hasMatch(symbol));
    _function = "refresh";

    _symbol = symbol;
    _$symbol.add(symbol);

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
      {String book = "1"}) {
    //assert(new RegExp(r"^[a-zA-Z]{2,6}(?:\d{2})*$").hasMatch(symbol));
    assert(new RegExp(r"^[a-zA-Z]{1,6}(?:\d{1,2})*$").hasMatch(symbol));
    assert(new RegExp(r"h|d|w|m").hasMatch(freq));
    //assert(new RegExp(r"^\d{8}$").hasMatch(time));
    assert(new RegExp(r"^(?:\d{4}|\d{8})$").hasMatch(time));

    _function = "refresh";

    _symbol = symbol;
    _frequency = freq;
    _time = time;

    if (_showRecords) {
      assert(new RegExp(r"^[a-zA-Z_0-9]+$").hasMatch(book));
      _book = book;
    }

    getChart();
  }

  void randomTradeRequest() {
    _function = "randomTrade";
    getChart();
  }

  void randomDateRequest() {
    _function = "randomDate";
    getChart();
  }

  String _requestUrl() {
    //var url = "http://127.0.0.1:5000/service/plot/practice";
    var url = "${window.location.origin}/service/plot/practice";

    url = "${url}?timestemp=${new DateTime.now().millisecondsSinceEpoch}";

    url =
        "${url}&symbol=${_symbol}&frequency=${_frequency}&function=${_function}&time=${_time}";

    //url = "${url}&book=${_book}";

    if (_showRecords) {
      url = "${url}&book=${_book}&records=true";
      //url = "${url}&records=true";
    }

    return url;
  }

  void infoRequest() async {
    if (_working) {
      return;
    }

    _function = "info";
    var url = _requestUrl();

    var info = await HttpRequest.getString(url);
    var m = json.decode(info);

    _$info.add(m);
    _$time.add(m["Time"]);
    _time = m["Time"];
  }

  void inspectRequest(num x, num y, {num ax, num ay}) async {
    assert(x >= 0 && y >= 0);

    if (_working) {
      return;
    }

    noteRequest(x, y);

    _function = "inspect";
    var url = _requestUrl();

    url = "${url}&x=${x}&y=${y}";

    if (ax != null && ay != null) {
      url = "${url}&ax=${ax}&ay=${ay}";
    }

    var info = await HttpRequest.getString(url);

    _$chartInspect.add(info);
  }

  void noteRequest(num x, num y) async {
    if (_working) {
      return;
    }

    var url = "${window.location.origin}/service/record/note";
    url = "${url}?timestemp=${new DateTime.now().millisecondsSinceEpoch}";
    url = "${url}&book=${_book}";
    url = "${url}&x=${x}&y=${y}";

    var note = await HttpRequest.getString(url);
    _note = note;
    _$note.add(note);
  }

  void getChart() {
    if (_working) {
      return;
    }

    _working = true;

    var url = _requestUrl();

    _$chartUrl.add(url);
  }

  bool get isWorking => _working;

  void done() {
    _working = false;
  }
}
