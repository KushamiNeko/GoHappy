import "../../_components/chart_inputs/chart_inputs.dart";
import "../../_components/canvas/canvas.dart";
import "../../_components/modal/modal.dart";
import "../../_components/navbar/navbar.dart";
import "../../_components/sidebar/sidebar.dart";
import "../../_services/control.dart";

void main() {
  final inputs = new ChartInputs("view");
  final sidebar = new Sidebar(inputs, "view");
  final modal = new Modal(inputs.container, "view");

  final navbar = new Navbar("view");

  new Canvas("view");

  MainControl(navbar, sidebar, modal, inputs);
}
