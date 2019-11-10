import "dart:html";
import "../../_services/server.dart";

class Note {
  final DivElement _container;
  final DivElement _content;

  bool _isOpen = false;

  final String _cls;

  final Server _server;

  Note(String id, {String cls = ""})
      : _cls = cls,
        _server = new Server(),
        _container = querySelector("#${id}-note"),
        _content = querySelector("#${id}-note-content") {
    close();

    _server.$note.listen((note) {
      _content.innerHtml = note;

      if (!_isOpen) {
        _container.style.bottom = "${-_container.clientHeight}px";
      }
    });
  }

  bool get isOpen => _isOpen;

  void open() {
    _container.style.bottom = "0px";
    _isOpen = true;
  }

  void close() {
    _container.style.bottom = "${-_container.clientHeight}px";
    _isOpen = false;
  }
}
