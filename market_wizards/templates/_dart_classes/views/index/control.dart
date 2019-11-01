import "dart:html";
import "../../components/chart_canvas/canvas.dart";
import "../../components/modal/modal.dart";
import "../../components/navbar/navbar.dart";
import "../../views/index/input.dart";
import "../../views/index/server.dart";

class IndexControl {
  final IndexServer _server;
  final IndexInput _inputs;

  final Navbar _navbar;
  final ChartCanvas _canvas;
  final Modal _modal;

  bool _isFullScreen = false;

  IndexControl()
      : _server = new IndexServer(),
        _inputs = new IndexInput(),
        _navbar = new Navbar(),
        _canvas = new ChartCanvas(),
        _modal = new Modal() {
    _modal.child = _inputs.container;

    window.onKeyDown.listen((KeyboardEvent event) {
      if (_inputs.isFocus()) {
        return;
      }

      _keyDownSwitch(event);
    });
    window.onKeyPress.listen((KeyboardEvent event) {
      if (_inputs.isFocus()) {
        return;
      }
      _keyPressSwitch(event);
    });

    _inputs.activate((String value) {
      _canvas.src = _server.fromInputs(value);
      //_canvas.info = _server.info();
      _inputs.setInputs(_server.info());
      _modal.close();
      return;
    });

    _canvas.image.onLoad.listen((Event event) async {
      _canvas.info = await _server.updateInfo();
      _inputs.setInputs(_server.info());
    });

    _navbar.activate((bool ans) {
      _server.showRecords = ans;
      _canvas.src = _server.imageUrl();
      //_canvas.info = _server.info();

      _inputs.showRecords = ans;
      _inputs.setInputs(_server.info());
    });

    _canvas.src = _server.imageUrl();
    //_canvas.info = _server.info();
    _inputs.setInputs(_server.info());
    _inputs.symbol = _server.currentSymbol;
  }

  void toggleFullScreen() {
    if (_isFullScreen) {
      querySelector("#index-navbar").classes.remove("index-modal-open");
      querySelector("#index-sidebar").classes.remove("index-modal-open");
      document.exitFullscreen();
    } else {
      querySelector("#index-navbar").classes.add("index-modal-open");
      querySelector("#index-sidebar").classes.add("index-modal-open");
      document.body.requestFullscreen();
    }
    _isFullScreen = !_isFullScreen;
  }

  void _keyDownSwitch(KeyboardEvent event) {
    switch (event.which) {
      case (27):
        // esc
        break;
      case (37):
        // left
        _canvas.src = _server.backward();
        break;
      case (39):
        // right
        _canvas.src = _server.forward();
        break;
    }

    _inputs.setInputs(_server.info());
    //_canvas.info = _server.info();
  }

  void _keyPressSwitch(KeyboardEvent event) {
    if (event.which >= 49 && event.which <= 57) {
      var i = event.which - 49;

      var symbols = _inputs.symbol.split(",");

      if (i >= 0 && i < symbols.length) {
        _canvas.src = _server.symbol(symbols[i]);
      }
      //0-9 number keys
    } else {
      switch (event.which) {
        case (104):
          // h
          _canvas.src = _server.frequency("h");
          break;
        case (100):
          // d
          _canvas.src = _server.frequency("d");
          break;
        case (119):
          // w
          _canvas.src = _server.frequency("w");
          break;
        case (109):
          // m
          _canvas.src = _server.frequency("m");
          break;
        case (13):
          // enter
          toggleFullScreen();
          break;
        case (32):
          // space
          if (_modal.isOpen) {
            _modal.close();
          } else {
            _modal.open();
          }
          break;
        default:
          break;
      }
    }

    _inputs.setInputs(_server.info());
    //_canvas.info = _server.info();
  }
}
