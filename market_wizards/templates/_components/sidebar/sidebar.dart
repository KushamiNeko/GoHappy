import "dart:html";
import "../chart_inputs/chart_inputs.dart";
//import "../../_services/server.dart";

class Sidebar {
  final DivElement _container;
  //final DivElement _note;

  final String _cls;

  //final Server _server;

  Sidebar(ChartInputs inputs, String id, {String cls = ""})
      : _cls = cls,
        //_server = new Server(),
        _container = querySelector("#${id}-sidebar-container") {}
  //_note = querySelector("#${id}-sidebar-note") {}

  void enterFullScreen(bool ans) {
    if (ans) {
      _container.classes.add("${_cls}-sidebar-hidden");
    } else {
      _container.classes.remove("${_cls}-sidebar-hidden");
    }
  }
}
