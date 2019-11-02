import "dart:html";

class Modal {
  final DivElement _modal;
  final DivElement _content;

  Element _parentElement;
  Element _childElement;

  bool _isOpen = false;

  final String _cls;

  Modal(Element content, String id, {String cls = ""})
      : _cls = cls,
        _modal = querySelector("#${id}-modal"),
        _content = querySelector("#${id}-modal-content") {
    _childElement = content;
    _parentElement = _childElement.parent;
  }

  bool get isOpen => _isOpen;

  void open() {
    _modal.classes.add("${_cls}-modal-open");
    _content.classes.add("${_cls}-modal-content-open");
    _childElement.remove();
    _content.children.add(_childElement);
    _isOpen = true;
  }

  void close() {
    _modal.classes.remove("${_cls}-modal-open");
    _content.classes.remove("${_cls}-modal-content-open");
    _childElement.remove();
    _parentElement.children.add(_childElement);
    _isOpen = false;
  }
}
