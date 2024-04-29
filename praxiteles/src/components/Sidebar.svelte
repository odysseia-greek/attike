<!-- Sidebar.svelte -->
<script>
    import Fa from 'svelte-fa'
    import {faBinoculars, faChartSimple} from '@fortawesome/free-solid-svg-icons'
    import { faAreaChart } from '@fortawesome/free-solid-svg-icons';
    import { faHomeUser, faChevronLeft, faChevronRight } from '@fortawesome/free-solid-svg-icons';

    let isSidebarOpen = false; // Set it as collapsed by default

    function toggleSidebar() {
        isSidebarOpen = !isSidebarOpen;
    }
</script>

<style>
    .sidebar {
        position: fixed;
        top: 0;
        left: 0; /* Display the sidebar on the left by default */
        width: 3.5em; /* Smaller width when collapsed */
        height: 100%;
        background-color: #333;
        z-index: 2;
        overflow-x: hidden;
        transition: width 0.3s ease-in-out; /* Add transition for width */
    }

    /* Add a border when collapsed */
    .collapsed {
        border-right: 1px solid #555;
    }

    /* Sidebar menu items */
    .menu-item {
        display: flex;
        align-items: center;
        padding: 1em;
        color: #e8e7e7;
        text-decoration: none;
        transition: background-color 0.2s;
    }

    .menu-item:hover {
        background-color: #555; /* Highlight on hover */
    }

    /* Icon styles (Font Awesome icons used for icons) */
    .menu-icon {
        margin-right: 1em;
        font-size: 1.5em;
    }

    /* Text to show when expanded */
    .text-expanded {
        display: block;
    }

    /* Icon to hide when collapsed */
    .icon-collapsed {
        display: none;
    }

    /* Adjust the sidebar width when expanded */
    .expanded {
        width: 15em; /* Slightly smaller width when expanded */
    }
</style>

<div class="sidebar" class:expanded={isSidebarOpen} class:collapsed={!isSidebarOpen}>
    <!-- Menu items with icons and text -->
    <a href="#top" class="menu-item">
        <span class="menu-icon"><Fa icon={faHomeUser} /></span>
        <span class={isSidebarOpen ? "text-expanded" : "icon-collapsed"}>Praxiteles</span>
    </a>
    <a href="#traces" class="menu-item">
        <span class="menu-icon"><Fa icon={faBinoculars} /></span>
        <span class={isSidebarOpen ? "text-expanded" : "icon-collapsed"}>Traces</span>
    </a>
    <a href="#traces-visual" class="menu-item">
        <span class="menu-icon"><Fa icon={faAreaChart} /></span>
        <span class={isSidebarOpen ? "text-expanded" : "icon-collapsed"}>Metrics</span>
    </a>
    <a href="#metrics" class="menu-item">
        <span class="menu-icon"><Fa icon={faChartSimple} /></span>
        <span class={isSidebarOpen ? "text-expanded" : "icon-collapsed"}>Metrics</span>
    </a>
    <a href="#metrics-visual" class="menu-item">
        <span class="menu-icon"><Fa icon={faAreaChart} /></span>
        <span class={isSidebarOpen ? "text-expanded" : "icon-collapsed"}>Metrics</span>
    </a>
    <!-- Add more menu items here as needed -->

    <!-- Icon to toggle expand/collapse within the sidebar -->
    <a class="menu-item" on:click={toggleSidebar}>
        <span class="menu-icon"><Fa icon={isSidebarOpen ? faChevronLeft : faChevronRight} /></span>
    </a>
</div>
