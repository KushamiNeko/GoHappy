import "dart:html";
import "../_components/chart_inputs/chart_inputs.dart";
import "../_components/trade_inputs/trade_inputs.dart";
import "../_components/modal/modal.dart";
import "../_components/canvas/canvas.dart";
import "../_components/note/note.dart";
import "../_components/navbar/navbar.dart";
import "../_components/sidebar/sidebar.dart";
import "server.dart";

class MainControl {
  final Server _server;
  final ChartInputs _chartInputs;
  final TradeInputs _tradeInputs;

  final Canvas _canvas;

  final Navbar _navbar;
  final Sidebar _sidebar;
  final Modal _modal;

  final Note _note;

  bool _isFullScreen = false;

  MainControl(this._navbar, this._sidebar, this._modal, this._canvas,
      this._chartInputs, this._note, this._tradeInputs)
      : _server = new Server() {
    window.onKeyDown.listen((KeyboardEvent event) {
      if (_chartInputs.isFocus()) {
        return;
      }

      _keyDownSwitch(event);
    });
    window.onKeyPress.listen((KeyboardEvent event) {
      if (_chartInputs.isFocus()) {
        return;
      }

      _keyPressSwitch(event);
    });

    document.body.onMouseDown.listen((MouseEvent event) {
      _canvas.mouseDown(event);
    });

    document.body.onMouseMove.listen((MouseEvent event) {
      _canvas.mouseMove(event);
      _tradeInputs.move(event.client.x, event.client.y);

      _server.noteRequest(event.client.x, event.client.y);
    });

    document.body.onMouseUp.listen((MouseEvent event) {
      _canvas.mouseUp(event);
    });
  }

  void toggleFullScreen() {
    _isFullScreen = !_isFullScreen;
    if (_isFullScreen) {
      document.body.requestFullscreen();
    } else {
      document.exitFullscreen();
    }

    _navbar.enterFullScreen(_isFullScreen);
    _sidebar.enterFullScreen(_isFullScreen);
  }

  void _keyDownSwitch(KeyboardEvent event) {
    switch (event.which) {
      case (27):
        // esc
        break;
      case (37):
        // left
        _server.backward();
        break;
      case (38):
        // up
        _chartInputs.symbolStep("b");
        break;
      case (39):
        // right
        _server.forward();
        break;
      case (40):
        // down
        _chartInputs.symbolStep("f");
        break;
    }
  }

  void _keyPressSwitch(KeyboardEvent event) {
    if (event.which >= 49 && event.which <= 57) {
      //1-9 number keys
      _chartInputs.symbolIndex(event.which - 49);
    } else {
      switch (event.which) {
        case (104):
          // h
          _server.freqRequest("h");
          break;
        case (100):
          // d
          _server.freqRequest("d");
          break;
        case (119):
          // w
          _server.freqRequest("w");
          break;
        case (109):
          // m
          _server.freqRequest("m");
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
        case (110):
          // n
          if (_note.isOpen) {
            _note.close();
          } else {
            _note.open();
          }
          break;

        case (116):
          // t
          if (_tradeInputs.isOpen()) {
            _tradeInputs.hide();
          } else {
            _tradeInputs.show();
          }
          break;
        default:
          break;
      }
    }
  }
}
