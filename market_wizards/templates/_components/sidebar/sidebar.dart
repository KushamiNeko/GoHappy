import "dart:html";
import "../chart_inputs/chart_inputs.dart";

class Sidebar {
  final DivElement _container;

  final String _cls;

  Sidebar(ChartInputs inputs, String id, {String cls = ""})
      : _cls = cls,
        _container = querySelector("#${id}-sidebar-container") {}

  void enterFullScreen(bool ans) {
    if (ans) {
      _container.classes.add("${_cls}-sidebar-hidden");
    } else {
      _container.classes.remove("${_cls}-sidebar-hidden");
    }
  }
}
