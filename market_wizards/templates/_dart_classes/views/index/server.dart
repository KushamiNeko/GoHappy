import "dart:html";

class IndexServer {
  String _time;
  String _freq = "d";
  String _symbol = "esz19";

  String _version = "1";

  String _action = "practice";
  String _function = "time";

  bool _showRecords = false;

  IndexServer() {
    var now = new DateTime.now();
    _time =
        "${now.year.toString()}${now.month.toString().padLeft(2, "0")}${now.day.toString().padLeft(2, "0")}"
            .padRight(14, "0");

    var path = window.location.pathname;

    if (path.contains("study")) {
      _action = "study";
    } else if (path.contains("practice")) {
      _action = "practice";
    }
  }

  void toggleShowRecords() {
    _showRecords = !_showRecords;
  }

  String forward() {
    _function = "forward";
    return imageUrl();
  }

  String backward() {
    _function = "backward";
    return imageUrl();
  }

  String time(String time) {
    _function = "time";
    var regex = new RegExp(r"\d{14}|\d{8}");

    if (!regex.hasMatch(time)) {
      throw Exception("invaid year ${time}");
    }

    var t = DateTime.parse(time.substring(0, 8));
    var now = new DateTime.now();

    if (t.isAfter(now)) {
      time =
          "${now.year}${now.month.toString().padLeft(2, "0")}${now.day.toString().padLeft(2, "0")}";
    }

    if (time.length == 8) {
      time = "${time}000000";
    }

    _time = time;
    return imageUrl();
  }

  String symbol(String symbol) {
    var regex = new RegExp(r"[a-zA-Z0-9]+");

    if (!regex.hasMatch(symbol)) {
      throw Exception("invaid symbol ${symbol}");
    }

    _symbol = symbol;
    return imageUrl();
  }

  String frequency(String freq) {
    _function = "frequency";
    var regex = new RegExp(r"h|d|w|m");

    if (!regex.hasMatch(freq)) {
      throw Exception("invaid frequency ${freq}");
    }

    _freq = freq;

    return imageUrl();
  }

  String version(String version) {
    if (_showRecords) {
      _version = version;
    }

    return imageUrl();
  }

  String fromInputs(String value) {
    var regex = new RegExp(r"\s+");
    value = value.replaceAll(regex, " ");

    var vs = value.split(" ");

    if (vs.length != 3 && vs.length != 4) {
      throw new Exception("invalid value: ${info}");
    }

    frequency(vs[2].trim().toLowerCase()[0]);
    symbol(vs[1].trim().toLowerCase());
    time(vs[0].trim());

    _function = "time";

    if (_showRecords) {
      version(vs[3].trim());
    }

    return imageUrl();
  }

  String _freqTable(String freq) {
    switch (freq) {
      case "h":
        return "hourly";
      case "d":
        return "daily";
      case "w":
        return "weekly";
      case "m":
        return "monthly";
      default:
        throw new Exception("unknown frequency: ${freq}");
    }
  }

  String info() {
    var time = _freq == "h" ? _time : _time.substring(0, 8);

    var info =
        "${time}  ${_symbol.toUpperCase()}  ${_freqTable(_freq).toUpperCase()}";

    if (_showRecords) {
      info = "${info}  ${_version}";
    }

    return info;
  }

  Future<String> updateInfo() async {
    _function = "info";
    var info = await HttpRequest.getString(imageUrl());
    var dtime =
        "${info.trim().split(" ")[0]}${info.trim().split(" ")[1].replaceAll(":", "")}";

    time(dtime);

    return info;
  }

  String imageUrl() {
    var time = _time.replaceAll(":", "").replaceAll(" ", "");
    var url = "/plot/${_action}/${_symbol}/${_freq}/${_function}/${time}";

    if (_showRecords) {
      url = "${url}/records/${_version}";
    }

    return url;
  }

  String get currentAction => _action;

  String get currentDate => _time;

  String get currentSymbol => _symbol;

  String get currentFrequency => _freq;

  String get currentVersion => _version;

  bool get showRecords => _showRecords;

  void set showRecords(bool ans) {
    _showRecords = ans;
  }
}
