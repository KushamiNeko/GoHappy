import "dart:html";

class Modal {
  final DivElement _modal = querySelector("#modal");
  final DivElement _content = querySelector("#modal-content");

  Element _parent;
  Element _child;

  bool _isOpen = false;

  Modal() {}

  void set child(Element child) {
    _child = child;
    _parent = _child.parent;
  }

  void open() {
    _modal.classes.add("modal-open");
    _content.classes.add("modal-content-open");
    _child.remove();
    _content.children.add(_child);
    _isOpen = true;
  }

  void close() {
    _modal.classes.remove("modal-open");
    _content.classes.remove("modal-content-open");
    _child.remove();
    _parent.children.add(_child);
    _isOpen = false;
  }

  bool get isOpen => _isOpen;
}
