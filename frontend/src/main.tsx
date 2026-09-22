import { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { routes } from "./router/routes";
import { DashboardPage } from "./pages/DashboardPage";
import { TurnaroundsPage } from "./pages/TurnaroundsPage";
import { TasksPage } from "./pages/TasksPage";
import { ResourcesPage } from "./pages/ResourcesPage";
import { DelaysPage } from "./pages/DelaysPage";
import "./styles.css";

const PAGES: Record<string, () => JSX.Element> = {
  "/dashboard": DashboardPage,
  "/turnarounds": TurnaroundsPage,
  "/tasks": TasksPage,
  "/resources": ResourcesPage,
  "/delays": DelaysPage
};

function App() {
  const [active, setActive] = useState<string>(routes[0]?.route ?? "/dashboard");
  const CurrentPage = PAGES[active] ?? DashboardPage;
  useEffect(() => {
    document.title = "航空地勤周转保障平台 · ground-turn";
  }, []);
  return (
    <div className="shell">
      <aside>
        <div className="brand">航空地勤周转<br />保障平台</div>
        <nav>
          {routes.map((route) => (
            <button
              key={route.route}
              className={active === route.route ? "active" : ""}
              onClick={() => setActive(route.route)}
            >
              {route.name}
            </button>
          ))}
        </nav>
      </aside>
      <main className="content">
        <CurrentPage />
      </main>
    </div>
  );
}

createRoot(document.getElementById("root")!).render(<App />);
