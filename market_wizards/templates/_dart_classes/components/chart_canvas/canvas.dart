import "dart:html";

class ChartCanvas {
  // final Element _container = querySelector("#container");
  final ImageElement _image = querySelector("#image");
  final SpanElement _info = querySelector("#info");

  final CanvasElement _inspectCanvas = querySelector("#inspect-canvas");
  final CanvasElement _coverCanvas = querySelector("#cover-canvas");

  CanvasRenderingContext2D _ictx;
  CanvasRenderingContext2D _cctx;

  final String _coverColor = "rgba(0, 0, 0, 0.8)";
  final String _inspectColor = "rgba(255, 255, 255, 0.8)";

  bool _cover = false;
  bool _double = false;
  num _doubleAnchorX = 0;

  ChartCanvas() {
    _ictx = _inspectCanvas.getContext("2d");
    _cctx = _coverCanvas.getContext("2d");

    _attachListener();
  }

  ImageElement get image => _image;

  void set src(String url) {
    _image.src = url;
    _image.onLoad.listen((Event event) {
      initCanvasSize();
    });
  }

  void set info(String info) {
    _info.innerHtml = info;
  }

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

    document.body.onMouseMove.listen((MouseEvent event) {
      if (_cover) {
        if (_double) {
          _doubleCover(event);
        } else {
          _singleCover(event);
        }
      } else {
        _inspect(event);
      }
    });

    document.body.onMouseDown.listen((MouseEvent event) {
      _ictx.clearRect(0, 0, _inspectCanvas.width, _inspectCanvas.height);
      _cctx.clearRect(0, 0, _inspectCanvas.width, _inspectCanvas.height);

      if (event.ctrlKey) {
        _singleCover(event);
      }

      if (event.altKey) {
        _doubleCover(event);
      }
    });

    _coverCanvas.onMouseUp.listen((MouseEvent event) {
      _cover = false;
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

  void _singleCover(MouseEvent event) {
    _cover = true;
    _cctx.clearRect(0, 0, _coverCanvas.width, _coverCanvas.height);
    _cctx.fillStyle = _coverColor;

    _cctx.fillRect(_eventXOffset(event), 0,
        _coverCanvas.width - _eventXOffset(event), _coverCanvas.height);
  }

  void _doubleCover(MouseEvent event) {
    _cctx.clearRect(0, 0, _coverCanvas.width, _coverCanvas.height);
    _cctx.fillStyle = _coverColor;

    if (!_double) {
      _doubleAnchorX = _eventXOffset(event);
    }

    _cover = true;
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
}
