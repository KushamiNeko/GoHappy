import "../../_components/chart_inputs/chart_inputs.dart";
import "../../_components/trade_inputs/trade_inputs.dart";
import "../../_components/canvas/canvas.dart";
import "../../_components/modal/modal.dart";
import "../../_components/note/note.dart";
import "../../_components/navbar/navbar.dart";
import "../../_components/sidebar/sidebar.dart";
import "../../_services/control.dart";

void main() {
  final canvas = new Canvas("view");

  final note = new Note("view");

  final cinputs = new ChartInputs("view");
  final sidebar = new Sidebar(cinputs, "view");
  final modal = new Modal(cinputs.container, "view");

  final tinputs = new TradeInputs("view");

  final navbar = new Navbar("view");

  MainControl(navbar, sidebar, modal, canvas, cinputs, note, tinputs);
}
