import "dart:html";
import "dart:math";
import "../../_services/server.dart";

class Canvas {
  final DivElement _info;

  //final DivElement _container;
  final ImageElement _image;

  final CanvasElement _inspectCanvas;
  final CanvasElement _coverCanvas;

  final Server _server;

  final String _cls;

  CanvasRenderingContext2D _ictx;
  CanvasRenderingContext2D _cctx;

  bool _calculate = false;

  num _calcAnchorX = 0;
  num _calcAnchorY = 0;

  bool _coverL = false;
  bool _coverR = false;

  bool _double = false;
  num _doubleAnchorX = 0;

  Canvas(String id, {String cls = ""})
      : _cls = cls,
        _server = new Server(),
        _info = querySelector("#${id}-canvas-chart-info"),
        //_container = querySelector("#${id}-canvas-chart-container"),
        _image = querySelector("#${id}-canvas-chart-image"),
        _inspectCanvas = querySelector("#${id}-canvas-chart-inspect"),
        _coverCanvas = querySelector("#${id}-canvas-chart-cover") {
    _ictx = _inspectCanvas.getContext("2d");
    _cctx = _coverCanvas.getContext("2d");

    _attachListener();
  }

  String get _coverColor => "rgba(0, 0, 0, 0.7)";
  String get _inspectColor => "rgba(255, 255, 255, 0.8)";
  String get _anchorColor => "rgba(255, 255, 255, 0.5)";

  void initCanvasSize() {
    _inspectCanvas.width = _image.clientWidth.floor() - 1;
    _inspectCanvas.height = _image.clientHeight.floor() - 1;

    _coverCanvas.width = _image.clientWidth.floor() - 1;
    _coverCanvas.height = _image.clientHeight.floor() - 1;
  }

  void _attachListener() {
    window.onResize.listen((Event event) {
      initCanvasSize();
    });

    _image.onLoad.listen((Event event) {
      //_container.classes.remove("${_cls}-canvas-chart-hidden");
      initCanvasSize();

      _server.done();

      _server.infoRequest();

      if (!_info.classes.contains("${_cls}-canvas-chart-info-hidden")) {
        _info.classes.add("${_cls}-canvas-chart-info-hidden");
      }
    });

    _server.$chartUrl.listen((url) {
      //_container.classes.add("${_cls}-canvas-chart-hidden");
      _image.src = url;
    });

    _server.done();
    _server.getChart();

    _server.$chartInspect.listen((info) {
      _info.innerHtml = info;
    });

    document.body.onMouseMove.listen((MouseEvent event) {
      _inspectInfo(event);

      if (_double) {
        _doubleCover(event);
      } else if (_coverL) {
        _singleCoverL(event);
      } else if (_coverR) {
        _singleCoverR(event);
      } else {}
      _inspect(event);
      if (_calculate) {
        _calcAnchor(event);
      }
    });

    document.body.onMouseDown.listen((MouseEvent event) {
      _ictx.clearRect(0, 0, _inspectCanvas.width, _inspectCanvas.height);
      _cctx.clearRect(0, 0, _inspectCanvas.width, _inspectCanvas.height);

      if (!_info.classes.contains("${_cls}-canvas-chart-info-hidden")) {
        _info.classes.add("${_cls}-canvas-chart-info-hidden");
      }

      _calculate = true;
      _calcAnchorX = _eventXOffset(event);
      _calcAnchorY = _eventYOffset(event);

      if (event.ctrlKey) {
        _singleCoverL(event);
      }

      if (event.shiftKey) {
        _singleCoverR(event);
      }

      if (event.altKey) {
        _doubleCover(event);
      }
    });

    _coverCanvas.onMouseUp.listen((MouseEvent event) {
      _calculate = false;

      _coverL = false;
      _coverR = false;

      _double = false;
      _doubleAnchorX = 0;
    });
  }

  num _eventXOffset(MouseEvent event) {
    return event.client.x - _image.offsetLeft;
  }

  num _eventYOffset(MouseEvent event) {
    return event.client.y - _image.offsetTop;
  }

  void _singleCoverR(MouseEvent event) {
    _coverR = true;

    _cctx.clearRect(0, 0, _coverCanvas.width, _coverCanvas.height);
    _cctx.fillStyle = _coverColor;

    _cctx.fillRect(_eventXOffset(event), 0,
        _coverCanvas.width - _eventXOffset(event), _coverCanvas.height);
  }

  void _singleCoverL(MouseEvent event) {
    _coverL = true;

    _cctx.clearRect(0, 0, _coverCanvas.width, _coverCanvas.height);
    _cctx.fillStyle = _coverColor;

    _cctx.fillRect(0, 0, _eventXOffset(event), _coverCanvas.height);
  }

  void _doubleCover(MouseEvent event) {
    _cctx.clearRect(0, 0, _coverCanvas.width, _coverCanvas.height);
    _cctx.fillStyle = _coverColor;

    if (!_double) {
      _doubleAnchorX = _eventXOffset(event);
    }

    _double = true;

    if (_eventXOffset(event) >= _doubleAnchorX) {
      _cctx.fillRect(0, 0, _doubleAnchorX, _coverCanvas.height);

      _cctx.fillRect(_eventXOffset(event), 0,
          _coverCanvas.width - _eventXOffset(event), _coverCanvas.height);
    } else {
      _cctx.fillRect(0, 0, _eventXOffset(event), _coverCanvas.height);

      _cctx.fillRect(_doubleAnchorX, 0, _coverCanvas.width - _doubleAnchorX,
          _coverCanvas.height);
    }
  }

  void _inspectInfo(MouseEvent event) {
    if (_info.classes.contains("${_cls}-canvas-chart-info-hidden")) {
      _info.classes.remove("${_cls}-canvas-chart-info-hidden");
    }

    var x = max(min(_eventXOffset(event) / _inspectCanvas.width, 1), 0);
    var y = max(
        min(
            (_inspectCanvas.height - _eventYOffset(event)) /
                _inspectCanvas.height,
            1),
        0);

    if (_calculate) {
      var ax = max(min(_calcAnchorX / _inspectCanvas.width, 1), 0);

      var ay = max(
          min((_inspectCanvas.height - _calcAnchorY) / _inspectCanvas.height,
              1),
          0);

      _server.inspectRequest(x, y, ax: ax, ay: ay);
    } else {
      _server.inspectRequest(x, y);
    }

    var offset = 20;

    if (_eventXOffset(event) > _inspectCanvas.width / 2) {
      _info.style.left = "${event.client.x - _info.clientWidth - offset}px";
    } else {
      _info.style.left = "${event.client.x + offset}px";
    }

    if (_eventYOffset(event) > _inspectCanvas.height / 2) {
      _info.style.top = "${event.client.y - _info.offsetHeight - offset}px";
    } else {
      _info.style.top = "${event.client.y + offset}px";
    }
  }

  void _inspect(MouseEvent event) {
    _ictx.clearRect(0, 0, _inspectCanvas.width, _inspectCanvas.height);
    _ictx.strokeStyle = _inspectColor;

    _ictx.beginPath();
    _ictx.moveTo(_eventXOffset(event), 0);

    _ictx.lineTo(_eventXOffset(event), _inspectCanvas.height);

    _ictx.stroke();
    _ictx.closePath();

    _ictx.beginPath();

    _ictx.moveTo(0, _eventYOffset(event));

    _ictx.lineTo(_inspectCanvas.width, _eventYOffset(event));

    _ictx.stroke();
    _ictx.closePath();
  }

  void _calcAnchor(MouseEvent event) {
    _ictx.strokeStyle = _anchorColor;

    _ictx.beginPath();
    _ictx.moveTo(_calcAnchorX, 0);

    _ictx.lineTo(_calcAnchorX, _inspectCanvas.height);

    _ictx.stroke();
    _ictx.closePath();

    _ictx.beginPath();

    _ictx.moveTo(0, _calcAnchorY);

    _ictx.lineTo(_inspectCanvas.width, _calcAnchorY);

    _ictx.stroke();
    _ictx.closePath();
  }
}
