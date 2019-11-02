import "dart:html";
import "../../_services/server.dart";

class ChartInputs {
  final Element _container;

  final InputElement _itime;
  final InputElement _ifreq;
  final InputElement _iversion;

  final TableSectionElement _tsymbols;
  final InputElement _isymbol;
  final ButtonElement _bsymbol;

  final ButtonElement _btn;

  final Server _server;

  final String _cls;

  TableCellElement _selectedCell;

  final List<String> _symbols = [
    //"spx",
    //"compq",
    //"rut",
    "esz19",
    "nqz19",
    "qrz19",
    "znz19",
    "clx19",
    "gcx19",
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
        _iversion = querySelector("#${id}-chart-inputs-version"),
        _btn = querySelector("#${id}-chart-inputs-button") {
    _server.$showRecords.listen((show) {
      if (show) {
        _iversion.parent.classes.remove("${cls}-chart-inputs-text-hidden");
      } else {
        _iversion.parent.classes.add("${cls}-chart-inputs-text-hidden");
      }
    });

    _server.$time.listen((time) {
      _itime.value = time;
    });

    _server.$frequency.listen((freq) {
      _ifreq.value = freq;
    });

    _server.$version.listen((verison) {
      _iversion.value = verison;
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
          version: _iversion.value);

      _btn.blur();
    });
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
        focused == _iversion ||
        focused == _btn);
  }
}
