import "dart:html";
import "../../_services/server.dart";

class TradeInputs {
  final DivElement _container;

  final InputElement _ibook;
  final InputElement _itime;
  final InputElement _isymbol;
  final InputElement _iprice;
  final InputElement _iquantity;

  final ButtonElement _binc;
  final ButtonElement _bdec;

  final ButtonElement _bok;

  final String _cls;

  final Server _server;

  TradeInputs(String id, {String cls = ""})
      : _cls = cls,
        _server = new Server(),
        _container = querySelector("#${id}-trade-inputs-container"),
        _ibook = querySelector("#${id}-trade-inputs-book"),
        _itime = querySelector("#${id}-trade-inputs-time"),
        _isymbol = querySelector("#${id}-trade-inputs-symbol"),
        _iprice = querySelector("#${id}-trade-inputs-price"),
        _iquantity = querySelector("#${id}-trade-inputs-quantity"),
        _binc = querySelector("#${id}-trade-inputs-increase"),
        _bdec = querySelector("#${id}-trade-inputs-decrease"),
        _bok = querySelector("#${id}-trade-inputs-ok") {
    _iquantity.value = "100";

    _binc.classes.add("${_cls}-trade-inputs-operation-btn-selected");
    _binc.onClick.listen((MouseEvent event) {
      _binc.classes.add("${_cls}-trade-inputs-operation-btn-selected");
      _bdec.classes.remove("${_cls}-trade-inputs-operation-btn-selected");
      _binc.blur();
    });

    _bdec.onClick.listen((MouseEvent event) {
      _bdec.classes.add("${_cls}-trade-inputs-operation-btn-selected");
      _binc.classes.remove("${_cls}-trade-inputs-operation-btn-selected");
      _bdec.blur();
    });

    _bok.onClick.listen((MouseEvent event) {
      hide();
    });

    _server.$book.listen((book) {
      _ibook.value = book;
    });

    _server.$symbol.listen((symbol) {
      var regex = new RegExp(r"^([a-z]+)(?:[fghjkmnquvxz])\d{2}$");
      var lsymbol = symbol.toLowerCase();

      if (regex.hasMatch(lsymbol)) {
        _isymbol.value = regex.firstMatch(lsymbol).group(1);
      }
    });

    _server.$time.listen((time) {
      _itime.value = time;
    });

    _server.$info.listen((info) {
      _iprice.value = info["Close"].toString();
    });
  }

  void move(num x, num y) {
    var width = _container.clientWidth;
    var height = _container.clientHeight;

    if (x >= document.body.clientWidth / 2.0) {
      _container.style.left = "${x - width}px";
    } else {
      _container.style.left = "${x}px";
    }

    if (y >= document.body.clientHeight / 2.0) {
      _container.style.top = "${y - height}px";
    } else {
      _container.style.top = "${y}px";
    }
  }

  void show() {
    _container.classes.remove("${_cls}-trade-inputs-hide");
  }

  void hide() {
    _container.classes.add("${_cls}-trade-inputs-hide");
  }

  bool isOpen() {
    return !_container.classes.contains("${_cls}-trade-inputs-hide");
  }
}
