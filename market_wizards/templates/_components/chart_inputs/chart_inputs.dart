import "dart:html";
import "../../_services/server.dart";

class ChartInputs {
  final Element _container;

  final InputElement _itime;
  final InputElement _ifreq;
  final InputElement _ibook;

  final TableSectionElement _tsymbols;
  final InputElement _isymbol;
  final ButtonElement _bsymbol;

  final ButtonElement _btn;

  //final ButtonElement _btnRD;
  final ButtonElement _btnRT;

  final Server _server;

  final String _cls;

  TableCellElement _selectedCell;

  final List<String> _symbols = [
    "es",
    "vix",
    "nq",
    "vxn",
    "qr",
    "rvx",
    "vle",
    "zn",
    "tyvix",
    "fx",
    "vstx",
    "np",
    "jniv",
    "cl",
    "ovx",
    "gc",
    "gvz",
    "mb",
    "vx",
    "rxes",
  ];

  ChartInputs(String id, {String cls = ""})
      : _cls = cls,
        _server = new Server(),
        _container = querySelector("#${id}-chart-inputs-container"),
        _tsymbols = querySelector("#${id}-chart-inputs-symbols-table"),
        _isymbol = querySelector("#${id}-chart-inputs-symbol"),
        _bsymbol = querySelector("#${id}-chart-inputs-symbol-button"),
        _itime = querySelector("#${id}-chart-inputs-time"),
        _ifreq = querySelector("#${id}-chart-inputs-frequency"),
        _ibook = querySelector("#${id}-chart-inputs-book"),
        _btn = querySelector("#${id}-chart-inputs-button"),
        //_btnRD = querySelector("#${id}-chart-inputs-random-date"),
        _btnRT = querySelector("#${id}-chart-inputs-random-trade") {
    _server.$showRecords.listen((show) {
      if (show) {
        _ibook.parent.classes.remove("${cls}-chart-inputs-text-hidden");
      } else {
        _ibook.parent.classes.add("${cls}-chart-inputs-text-hidden");
      }
    });

    _server.$time.listen((time) {
      _itime.value = time;
    });

    _server.$frequency.listen((freq) {
      _ifreq.value = freq;
    });

    _server.$book.listen((book) {
      _ibook.value = book;
    });

    _server.broadcast();

    for (var symbol in _symbols) {
      _addSymbol(symbol);
    }

    _bsymbol.onClick.listen((MouseEvent event) {
      _isymbol.classes.remove("${cls}-chart-inputs-text-error");

      var symbol = _isymbol.value;
      var regex = new RegExp(r"^[a-zA-Z]{2,6}(?:\d{2})*$");

      if (!regex.hasMatch(symbol)) {
        _isymbol.classes.add("${cls}-chart-inputs-text-error");
        return;
      } else {
        _addSymbol(symbol);
        _isymbol.value = "";
      }

      _bsymbol.blur();
    });

    assert(_tsymbols.children[0].children.length == 1);

    _tsymbols.children[0].children[0].click();

    _btn.onClick.listen((MouseEvent event) {
      _server.inputsRequest(_selectedCell.innerHtml, _itime.value, _ifreq.value,
          book: _ibook.value);

      _btn.blur();
    });

    _btnRT.onClick.listen((MouseEvent event) {
      _server.randomTradeRequest();
    });

    //_btnRD.onClick.listen((MouseEvent event) {
    //_server.randomDateRequest();
    //});
  }

  Element get container => _container;

  void _addSymbol(String symbol) {
    var tr = _tsymbols.addRow();
    var td = tr.addCell();

    td.innerHtml = symbol;

    td.onClick.listen((Event event) {
      var t = event.target as TableCellElement;
      _selectedCell?.classes?.remove("${_cls}-chart-inputs-table-selected");
      t.classes.add("${_cls}-chart-inputs-table-selected");
      _selectedCell = t;

      _server.symbolRequest(_selectedCell.innerHtml);
    });
  }

  void symbolIndex(int id) {
    if (_server.isWorking) {
      return;
    }

    if (id < _tsymbols.children.length) {
      _tsymbols.children[id].children[0].click();
    }
  }

  void symbolStep(String dir) {
    assert(dir == "f" || dir == "b");

    if (_server.isWorking) {
      return;
    }

    var id = _tsymbols.children.indexOf(_selectedCell.parent);

    id = dir == "f" ? id + 1 : id - 1;
    id = id % _tsymbols.children.length;

    assert(id >= 0 && id <= _tsymbols.children.length - 1);

    _tsymbols.children[id].children[0].click();
  }

  bool isFocus() {
    var focused = document.activeElement;
    return (focused == _itime ||
        focused == _tsymbols ||
        focused == _bsymbol ||
        focused == _isymbol ||
        focused == _ifreq ||
        focused == _ibook ||
        focused == _btn);
  }
}
