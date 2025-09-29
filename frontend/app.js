// Импортируем API Wails из сгенерированных bindings
import { StartMonitoring, StopMonitoring, ResetStats, GetStats, GetAbilities, GetTargets, OpenDevTools } from './wailsjs/wailsjs/go/app/App.js';

class DPSMeter {
    constructor() {
        this.isMonitoring = true;
        this.updateInterval = null;
        this.abilitiesSort = { column: 'damage', direction: 'desc' };
        this.targetsSort = { column: 'damage', direction: 'desc' };
        this.abilitiesCollapsed = false;
        this.targetsCollapsed = false;
        this.initializeElements();
        this.bindEvents();
        this.updateStatus('Ready to start monitoring');
    }

    initializeElements() {
        this.startBtn = document.getElementById('startBtn');
        this.stopBtn = document.getElementById('stopBtn');
        this.resetBtn = document.getElementById('resetBtn');
        this.debugBtn = document.getElementById('debugBtn');
        this.statusText = document.getElementById('statusText');
        this.debugPanel = document.getElementById('debugPanel');
        this.debugContent = document.getElementById('debugContent');
        
        // Stats elements
        this.maxDpsElement = document.getElementById('maxDps');
        this.dpsElement = document.getElementById('dps');
        this.damageElement = document.getElementById('damage');
        this.hitsElement = document.getElementById('hits');
        this.critsElement = document.getElementById('crits');
        this.maxHpsElement = document.getElementById('maxHps');
        this.hpsElement = document.getElementById('hps');
        this.healingElement = document.getElementById('healing');
        this.healingHitsElement = document.getElementById('healingHits');
        this.healingCritsElement = document.getElementById('healingCrits');
        this.critRateElement = document.getElementById('critRate');
        this.healingCritRateElement = document.getElementById('healingCritRate');
        this.killsElement = document.getElementById('kills');
        this.durationElement = document.getElementById('duration');
        
        // Table elements
        this.abilitiesTable = document.getElementById('abilitiesTable');
        this.targetsTable = document.getElementById('targetsTable');
        this.abilitiesContainer = document.getElementById('abilitiesContainer');
        this.targetsContainer = document.getElementById('targetsContainer');
        
        // Collapse buttons
        this.collapseAbilitiesBtn = document.getElementById('collapseAbilitiesBtn');
        this.collapseTargetsBtn = document.getElementById('collapseTargetsBtn');
    }

    bindEvents() {
        this.startBtn.addEventListener('click', () => this.startMonitoring());
        this.stopBtn.addEventListener('click', () => this.stopMonitoring());
        this.resetBtn.addEventListener('click', () => this.resetStats());
        this.debugBtn.addEventListener('click', () => this.showDebugInfo());
        
        // Добавляем обработчики сортировки для таблиц
        this.bindSortEvents();
        
        // Добавляем обработчики сворачивания
        this.bindCollapseEvents();
    }

    async startMonitoring() {
        try {
            console.log('Starting monitoring...');
            const result = await StartMonitoring();
            console.log('StartMonitoring result:', result);
            this.updateStatus(result);
            
            if (result.includes('started') || result.includes('monitoring')) {
                this.isMonitoring = true;
                this.startBtn.disabled = true;
                this.stopBtn.disabled = false;
                this.startUpdating();
            }
        } catch (error) {
            console.error('Error starting monitoring:', error);
            this.updateStatus('Error starting monitoring: ' + error.message);
        }
    }

    async stopMonitoring() {
        try {
            console.log('Stopping monitoring...');
            const result = await StopMonitoring();
            console.log('StopMonitoring result:', result);
            this.updateStatus(result);
            
            this.isMonitoring = false;
            this.startBtn.disabled = false;
            this.stopBtn.disabled = true;
            this.stopUpdating();
        } catch (error) {
            console.error('Error stopping monitoring:', error);
            this.updateStatus('Error stopping monitoring: ' + error.message);
        }
    }

    async resetStats() {
        try {
            const result = await ResetStats();
            this.updateStatus(result);
            await this.updateStats();
        } catch (error) {
            console.error('Error resetting stats:', error);
            this.updateStatus('Error resetting stats: ' + error.message);
        }
    }

    async showDebugInfo() {
        try {
            console.log('=== DEBUG INFO ===');
            
            const debugInfo = {
                isMonitoring: this.isMonitoring,
                updateInterval: this.updateInterval,
                timestamp: new Date().toISOString()
            };
            
            try {
                const stats = await GetStats();
                debugInfo.stats = stats;
                console.log('Current Stats:', stats);
            } catch (e) {
                debugInfo.statsError = e.message;
            }
            
            try {
                const abilities = await GetAbilities();
                debugInfo.abilities = abilities;
                console.log('Abilities:', abilities);
            } catch (e) {
                debugInfo.abilitiesError = e.message;
            }
            
            try {
                const targets = await GetTargets();
                debugInfo.targets = targets;
                console.log('Targets:', targets);
            } catch (e) {
                debugInfo.targetsError = e.message;
            }
            
            // Показываем информацию в отладочной панели
            this.debugContent.innerHTML = `<pre>${JSON.stringify(debugInfo, null, 2)}</pre>`;
            this.debugPanel.style.display = 'block';
            
            // Показываем информацию в статусе
            this.updateStatus(`Debug: Monitoring=${this.isMonitoring}, Stats loaded`);
            
            // Попытка открыть DevTools через Wails API
            try {
                const result = await OpenDevTools();
                console.log('OpenDevTools result:', result);
            } catch (e) {
                console.log('Could not open DevTools:', e);
            }
            
        } catch (error) {
            console.error('Error getting debug info:', error);
            this.updateStatus('Debug error: ' + error.message);
            this.debugContent.innerHTML = `<pre>Error: ${error.message}</pre>`;
            this.debugPanel.style.display = 'block';
        }
    }

    startUpdating() {
        // Обновляем статистику каждую секунду
        this.updateInterval = setInterval(() => {
            this.updateStats();
        }, 1000);
        
        // Первое обновление сразу
        this.updateStats();
    }

    stopUpdating() {
        if (this.updateInterval) {
            clearInterval(this.updateInterval);
            this.updateInterval = null;
        }
    }

    async updateStats() {
        try {
            const [stats, abilities, targets] = await Promise.all([
                GetStats(),
                GetAbilities(),
                GetTargets()
            ]);

            // Сохраняем данные для сортировки
            this.lastAbilities = abilities;
            this.lastTargets = targets;

            // console.log('Stats:', stats); // Раскомментируйте для отладки
            this.updateStatsDisplay(stats);
            this.updateAbilitiesTable(abilities);
            this.updateTargetsTable(targets);
        } catch (error) {
            console.error('Error updating stats:', error);
        }
    }

    updateStatsDisplay(stats) {
        // Обновляем значения
        this.maxDpsElement.textContent = this.formatNumber(stats.maxDps);
        this.dpsElement.textContent = this.formatNumber(stats.dps);
        this.damageElement.textContent = this.formatDamage(stats.damage);
        this.hitsElement.textContent = stats.hits;
        this.critsElement.textContent = stats.crits;
        this.maxHpsElement.textContent = this.formatNumber(stats.maxHps);
        this.hpsElement.textContent = this.formatNumber(stats.hps);
        this.healingElement.textContent = this.formatDamage(stats.healing);
        this.healingHitsElement.textContent = stats.healingHits;
        this.healingCritsElement.textContent = stats.healingCrits;
        this.critRateElement.textContent = this.formatPercentage(stats.critRate);
        this.healingCritRateElement.textContent = this.formatPercentage(stats.healingCritRate);
        this.killsElement.textContent = stats.kills;
        this.durationElement.textContent = this.formatDuration(stats.duration);
    }

    updateAbilitiesTable(abilities) {
        this.abilitiesTable.innerHTML = '';
        
        // Сортируем данные
        const sortedAbilities = this.sortData(abilities, this.abilitiesSort);
        
        sortedAbilities.forEach(ability => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${this.escapeHtml(ability.name)}</td>
                <td>${this.formatDamage(ability.damage)}</td>
                <td>${this.formatDamage(ability.healing)}</td>
                <td>${ability.hits}</td>
                <td>${ability.crits}</td>
                <td>${this.formatPercentage(ability.critRate || 0)}</td>
                <td>${ability.healingHits || 0}</td>
                <td>${ability.healingCrits || 0}</td>
                <td>${this.formatPercentage(ability.healingCritRate || 0)}</td>
            `;
            this.abilitiesTable.appendChild(row);
        });
    }

    updateTargetsTable(targets) {
        this.targetsTable.innerHTML = '';
        
        // Сортируем данные
        const sortedTargets = this.sortData(targets, this.targetsSort);
        
        sortedTargets.forEach(target => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${this.escapeHtml(target.name)}</td>
                <td>${this.formatDamage(target.damage)}</td>
                <td>${this.formatDamage(target.healing)}</td>
                <td>${target.hits}</td>
                <td>${target.crits}</td>
                <td>${this.formatPercentage(target.critRate || 0)}</td>
                <td>${target.kills}</td>
                <td>${target.healingHits || 0}</td>
                <td>${target.healingCrits || 0}</td>
                <td>${this.formatPercentage(target.healingCritRate || 0)}</td>
            `;
            this.targetsTable.appendChild(row);
        });
    }

    updateStatus(message) {
        this.statusText.textContent = message;
        console.log('Status:', message);
    }

    formatNumber(num) {
        if (typeof num !== 'number') return '0';
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + 'M';
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + 'K';
        }
        return Math.round(num).toLocaleString();
    }

    formatDamage(num) {
        if (typeof num !== 'number') return '0';
        // Для урона всегда показываем точные числа без округления
        return num.toLocaleString();
    }

    formatPercentage(num) {
        if (typeof num !== 'number') return '0%';
        return num.toFixed(1) + '%';
    }

    formatDuration(seconds) {
        if (typeof seconds !== 'number') return '0s';
        
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        const secs = Math.floor(seconds % 60);
        
        if (hours > 0) {
            return `${hours}h ${minutes}m ${secs}s`;
        } else if (minutes > 0) {
            return `${minutes}m ${secs}s`;
        } else {
            return `${secs}s`;
        }
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    // Методы для сортировки
    bindSortEvents() {
        // Обработчики для таблицы Abilities
        const abilitiesHeaders = this.abilitiesTable.parentElement.querySelectorAll('th.sortable');
        abilitiesHeaders.forEach(header => {
            header.addEventListener('click', () => {
                const column = header.dataset.column;
                this.handleSort('abilities', column);
            });
        });

        // Обработчики для таблицы Targets
        const targetsHeaders = this.targetsTable.parentElement.querySelectorAll('th.sortable');
        targetsHeaders.forEach(header => {
            header.addEventListener('click', () => {
                const column = header.dataset.column;
                this.handleSort('targets', column);
            });
        });
    }

    handleSort(tableType, column) {
        if (tableType === 'abilities') {
            if (this.abilitiesSort.column === column) {
                // Переключаем направление сортировки
                this.abilitiesSort.direction = this.abilitiesSort.direction === 'asc' ? 'desc' : 'asc';
            } else {
                // Новая колонка, начинаем с убывания
                this.abilitiesSort.column = column;
                this.abilitiesSort.direction = 'desc';
            }
            this.updateSortIndicators('abilities');
            this.updateAbilitiesTable(this.lastAbilities || []);
        } else if (tableType === 'targets') {
            if (this.targetsSort.column === column) {
                // Переключаем направление сортировки
                this.targetsSort.direction = this.targetsSort.direction === 'asc' ? 'desc' : 'asc';
            } else {
                // Новая колонка, начинаем с убывания
                this.targetsSort.column = column;
                this.targetsSort.direction = 'desc';
            }
            this.updateSortIndicators('targets');
            this.updateTargetsTable(this.lastTargets || []);
        }
    }

    updateSortIndicators(tableType) {
        const table = tableType === 'abilities' ? 
            this.abilitiesTable.parentElement : 
            this.targetsTable.parentElement;
        const sort = tableType === 'abilities' ? 
            this.abilitiesSort : 
            this.targetsSort;

        // Сбрасываем все индикаторы
        const headers = table.querySelectorAll('th.sortable');
        headers.forEach(header => {
            header.classList.remove('sorted-asc', 'sorted-desc');
            const indicator = header.querySelector('.sort-indicator');
            indicator.textContent = '↕';
        });

        // Устанавливаем индикатор для активной колонки
        const activeHeader = table.querySelector(`th[data-column="${sort.column}"]`);
        if (activeHeader) {
            const indicator = activeHeader.querySelector('.sort-indicator');
            if (sort.direction === 'asc') {
                activeHeader.classList.add('sorted-asc');
                indicator.textContent = '↑';
            } else {
                activeHeader.classList.add('sorted-desc');
                indicator.textContent = '↓';
            }
        }
    }

    sortData(data, sortConfig) {
        return [...data].sort((a, b) => {
            let aVal = a[sortConfig.column];
            let bVal = b[sortConfig.column];

            // Обрабатываем числовые значения
            if (typeof aVal === 'number' && typeof bVal === 'number') {
                return sortConfig.direction === 'asc' ? aVal - bVal : bVal - aVal;
            }

            // Обрабатываем строковые значения
            aVal = String(aVal || '').toLowerCase();
            bVal = String(bVal || '').toLowerCase();

            if (sortConfig.direction === 'asc') {
                return aVal.localeCompare(bVal);
            } else {
                return bVal.localeCompare(aVal);
            }
        });
    }

    // Методы для сворачивания
    bindCollapseEvents() {
        this.collapseAbilitiesBtn.addEventListener('click', () => {
            this.toggleCollapse('abilities');
        });

        this.collapseTargetsBtn.addEventListener('click', () => {
            this.toggleCollapse('targets');
        });
    }

    toggleCollapse(section) {
        if (section === 'abilities') {
            this.abilitiesCollapsed = !this.abilitiesCollapsed;
            this.updateCollapseState('abilities', this.abilitiesCollapsed);
        } else if (section === 'targets') {
            this.targetsCollapsed = !this.targetsCollapsed;
            this.updateCollapseState('targets', this.targetsCollapsed);
        }
    }

    updateCollapseState(section, collapsed) {
        if (section === 'abilities') {
            if (collapsed) {
                this.abilitiesContainer.classList.add('collapsed');
                this.collapseAbilitiesBtn.classList.add('collapsed');
            } else {
                this.abilitiesContainer.classList.remove('collapsed');
                this.collapseAbilitiesBtn.classList.remove('collapsed');
            }
        } else if (section === 'targets') {
            if (collapsed) {
                this.targetsContainer.classList.add('collapsed');
                this.collapseTargetsBtn.classList.add('collapsed');
            } else {
                this.targetsContainer.classList.remove('collapsed');
                this.collapseTargetsBtn.classList.remove('collapsed');
            }
        }
    }
}

// Инициализация приложения
document.addEventListener('DOMContentLoaded', () => {
    window.dpsMeter = new DPSMeter();
});

// Обработка ошибок
window.addEventListener('error', (event) => {
    console.error('Global error:', event.error);
});

window.addEventListener('unhandledrejection', (event) => {
    console.error('Unhandled promise rejection:', event.reason);
});
