// Language switching functionality
window.switchLanguage = function(lang) {
    try {
        // Update language buttons
        var zhBtn = document.getElementById('zh-btn');
        var enBtn = document.getElementById('en-btn');
        if (zhBtn) zhBtn.classList.remove('active');
        if (enBtn) enBtn.classList.remove('active');
        var langBtn = document.getElementById(lang + '-btn');
        if (langBtn) langBtn.classList.add('active');
        
        // Update page title
        const pageTitle = document.getElementById('page-title');
        if (pageTitle) {
            if (lang === 'zh') {
                pageTitle.textContent = 'Ceph 日志分析';
            } else {
                pageTitle.textContent = 'Ceph Log Analysis';
            }
        }
        
        // Update tab labels
        const tabs = document.getElementsByClassName('tab');
        for (let i = 0; i < tabs.length; i++) {
            tabs[i].textContent = tabs[i].getAttribute('data-' + lang);
        }
        
        // Update AIO panel
        updateAIOPanel(lang);
        
        // Update Repop panel
        updateRepopPanel(lang);
        
        // Update OSD panel
        updateOSDPanel(lang);
        
        // Update Transaction panel
        updateTransactionPanel(lang);
        
        // Update Metadata panel
        updateMetadataPanel(lang);
        
        // Update Client panel
        updateClientPanel(lang);
        
        // Update Client Detail panel
        updateClientDetailPanel(lang);
        
        // Update Dequeue panel
        updateDequeuePanel(lang);
        
        // Update OSD Op panel
        updateOSDOpPanel(lang);
    } catch (e) {
        console.error("Error in switchLanguage:", e);
    }
};

function updateAIOPanel(lang) {
    try {
        const panel = document.getElementById('aio');
        if (!panel) return;
        const elements = panel.querySelectorAll('[data-zh]');
        elements.forEach(el => {
            el.textContent = el.getAttribute('data-' + lang);
        });
    } catch (e) {
        console.error("Error in updateAIOPanel:", e);
    }
}

function updateRepopPanel(lang) {
    try {
        const panel = document.getElementById('repop');
        if (!panel) return;
        const elements = panel.querySelectorAll('[data-zh]');
        elements.forEach(el => {
            el.textContent = el.getAttribute('data-' + lang);
        });
    } catch (e) {
        console.error("Error in updateRepopPanel:", e);
    }
}

function updateOSDPanel(lang) {
    try {
        const panel = document.getElementById('osd');
        if (!panel) return;
        const elements = panel.querySelectorAll('[data-zh]');
        elements.forEach(el => {
            el.textContent = el.getAttribute('data-' + lang);
        });
    } catch (e) {
        console.error("Error in updateOSDPanel:", e);
    }
}

function updateTransactionPanel(lang) {
    try {
        const panel = document.getElementById('transaction');
        if (!panel) return;
        const elements = panel.querySelectorAll('[data-zh]');
        elements.forEach(el => {
            el.textContent = el.getAttribute('data-' + lang);
        });
    } catch (e) {
        console.error("Error in updateTransactionPanel:", e);
    }
}

function updateMetadataPanel(lang) {
    try {
        const panel = document.getElementById('metadata');
        if (!panel) return;
        const elements = panel.querySelectorAll('[data-zh]');
        elements.forEach(el => {
            el.textContent = el.getAttribute('data-' + lang);
        });
    } catch (e) {
        console.error("Error in updateMetadataPanel:", e);
    }
}

function updateClientPanel(lang) {
    try {
        const panel = document.getElementById('client');
        if (!panel) return;
        const elements = panel.querySelectorAll('[data-zh]');
        elements.forEach(el => {
            el.textContent = el.getAttribute('data-' + lang);
        });
    } catch (e) {
        console.error("Error in updateClientPanel:", e);
    }
}

function updateClientDetailPanel(lang) {
    try {
        const panel = document.getElementById('client-detail');
        if (!panel) return;
        const elements = panel.querySelectorAll('[data-zh]');
        elements.forEach(el => {
            el.textContent = el.getAttribute('data-' + lang);
        });
    } catch (e) {
        console.error("Error in updateClientDetailPanel:", e);
    }
}

function updateDequeuePanel(lang) {
    try {
        const panel = document.getElementById('dequeue');
        if (!panel) return;
        const elements = panel.querySelectorAll('[data-zh]');
        elements.forEach(el => {
            el.textContent = el.getAttribute('data-' + lang);
        });
    } catch (e) {
        console.error("Error in updateDequeuePanel:", e);
    }
}

function updateOSDOpPanel(lang) {
    try {
        const panel = document.getElementById('osdop');
        if (!panel) return;
        const elements = panel.querySelectorAll('[data-zh]');
        elements.forEach(el => {
            el.textContent = el.getAttribute('data-' + lang);
        });
    } catch (e) {
        console.error("Error in updateOSDOpPanel:", e);
    }
}

window.openTab = function(evt, tabName) {
    try {
        var i, tabcontent, tablinks;
        
        // Hide all tab content
        tabcontent = document.getElementsByClassName("panel");
        for (i = 0; i < tabcontent.length; i++) {
            tabcontent[i].classList.remove("active");
        }
        
        // Remove active class from all tabs
        tablinks = document.getElementsByClassName("tab");
        for (i = 0; i < tablinks.length; i++) {
            tablinks[i].classList.remove("active");
        }
        
        // Show the selected tab content and set active tab
        var tabContent = document.getElementById(tabName);
        if (tabContent) tabContent.classList.add("active");
        if (evt && evt.currentTarget) evt.currentTarget.classList.add("active");
    } catch (e) {
        console.error("Error in openTab:", e);
    }
};

// AIO Table Filter
window.filterAIOTable = function() {
    try {
        var startTime = document.getElementById("aio-start-time").value;
        var endTime = document.getElementById("aio-end-time").value;
        var minDuration = document.getElementById("aio-min-duration").value;
        var maxDuration = document.getElementById("aio-max-duration").value;
        var blockType = document.getElementById("aio-block-type").value;
        var minLength = document.getElementById("aio-min-length").value;
        var maxLength = document.getElementById("aio-max-length").value;
        var table = document.getElementById("aio-table");
        if (!table) return;
        var tr = table.getElementsByTagName("tr");
        
        for (var i = 1; i < tr.length; i++) {
            var tdStartTime = tr[i].getElementsByTagName("td")[0].textContent;
            var tdEndTime = tr[i].getElementsByTagName("td")[1].textContent;
            var tdDuration = parseFloat(tr[i].getElementsByTagName("td")[2].textContent);
            var tdLength = parseInt(tr[i].getElementsByTagName("td")[4].textContent);
            var tdBlockType = tr[i].getElementsByTagName("td")[5].textContent;
            
            var match = true;
            
            if (startTime) {
                var startDate = new Date(startTime.replace('T', ' '));
                var rowStartDate = new Date(tdStartTime);
                if (rowStartDate < startDate) match = false;
            }
            
            if (endTime) {
                var endDate = new Date(endTime.replace('T', ' '));
                var rowEndDate = new Date(tdEndTime);
                if (rowEndDate > endDate) match = false;
            }
            
            if (minDuration) {
                if (tdDuration < parseFloat(minDuration)) match = false;
            }
            
            if (maxDuration) {
                if (tdDuration > parseFloat(maxDuration)) match = false;
            }
            
            if (blockType) {
                if (tdBlockType !== blockType) match = false;
            }
            
            if (minLength) {
                if (tdLength < parseInt(minLength)) match = false;
            }
            
            if (maxLength) {
                if (tdLength > parseInt(maxLength)) match = false;
            }
            
            tr[i].style.display = match ? "" : "none";
        }

        window.initPagination("aio-table");
    } catch (e) {
        console.error("Error in filterAIOTable:", e);
    }
};

window.resetAIOFilter = function() {
    try {
        var startTimeInput = document.getElementById("aio-start-time");
        var endTimeInput = document.getElementById("aio-end-time");
        var minDurationInput = document.getElementById("aio-min-duration");
        var maxDurationInput = document.getElementById("aio-max-duration");
        var blockTypeInput = document.getElementById("aio-block-type");
        var minLengthInput = document.getElementById("aio-min-length");
        var maxLengthInput = document.getElementById("aio-max-length");
        
        if (startTimeInput) startTimeInput.value = "";
        if (endTimeInput) endTimeInput.value = "";
        if (minDurationInput) minDurationInput.value = "";
        if (maxDurationInput) maxDurationInput.value = "";
        if (blockTypeInput) blockTypeInput.value = "";
        if (minLengthInput) minLengthInput.value = "";
        if (maxLengthInput) maxLengthInput.value = "";
        
        filterAIOTable();
    } catch (e) {
        console.error("Error in resetAIOFilter:", e);
    }
};

// Repop Table Filter
window.filterRepopTable = function() {
    try {
        var startTime = document.getElementById("repop-start-time").value;
        var endTime = document.getElementById("repop-end-time").value;
        var minDuration = document.getElementById("repop-min-duration").value;
        var maxDuration = document.getElementById("repop-max-duration").value;
        var table = document.getElementById("repop-table");
        if (!table) return;
        var tr = table.getElementsByTagName("tr");
        
        for (var i = 1; i < tr.length; i++) {
            var tdStartTime = tr[i].getElementsByTagName("td")[0].textContent;
            var tdEndTime = tr[i].getElementsByTagName("td")[1].textContent;
            var tdDuration = parseFloat(tr[i].getElementsByTagName("td")[2].textContent);
            
            var match = true;
            
            if (startTime) {
                var startDate = new Date(startTime.replace('T', ' '));
                var rowStartDate = new Date(tdStartTime);
                if (rowStartDate < startDate) match = false;
            }
            
            if (endTime) {
                var endDate = new Date(endTime.replace('T', ' '));
                var rowEndDate = new Date(tdEndTime);
                if (rowEndDate > endDate) match = false;
            }
            
            if (minDuration) {
                if (tdDuration < parseFloat(minDuration)) match = false;
            }
            
            if (maxDuration) {
                if (tdDuration > parseFloat(maxDuration)) match = false;
            }
            
            tr[i].style.display = match ? "" : "none";
        }

        window.initPagination("repop-table");
    } catch (e) {
        console.error("Error in filterRepopTable:", e);
    }
};

window.resetRepopFilter = function() {
    try {
        var startTimeInput = document.getElementById("repop-start-time");
        var endTimeInput = document.getElementById("repop-end-time");
        var minDurationInput = document.getElementById("repop-min-duration");
        var maxDurationInput = document.getElementById("repop-max-duration");
        
        if (startTimeInput) startTimeInput.value = "";
        if (endTimeInput) endTimeInput.value = "";
        if (minDurationInput) minDurationInput.value = "";
        if (maxDurationInput) maxDurationInput.value = "";
        
        filterRepopTable();
    } catch (e) {
        console.error("Error in resetRepopFilter:", e);
    }
};

// OSD Table Filter
window.filterOSDTable = function() {
    try {
        var startTime = document.getElementById("osd-start-time").value;
        var endTime = document.getElementById("osd-end-time").value;
        var minLatency = document.getElementById("osd-min-latency").value;
        var maxLatency = document.getElementById("osd-max-latency").value;
        var table = document.getElementById("osd-table");
        if (!table) return;
        var tr = table.getElementsByTagName("tr");
        
        for (var i = 1; i < tr.length; i++) {
            var tdTime = tr[i].getElementsByTagName("td")[0].textContent;
            var tdLatency = parseFloat(tr[i].getElementsByTagName("td")[8].textContent);
            
            var match = true;
            
            if (startTime) {
                var startDate = new Date(startTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate < startDate) match = false;
            }
            
            if (endTime) {
                var endDate = new Date(endTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate > endDate) match = false;
            }
            
            if (minLatency) {
                if (tdLatency < parseFloat(minLatency)) match = false;
            }
            
            if (maxLatency) {
                if (tdLatency > parseFloat(maxLatency)) match = false;
            }
            
            tr[i].style.display = match ? "" : "none";
        }

        window.initPagination("osd-table");
    } catch (e) {
        console.error("Error in filterOSDTable:", e);
    }
};

window.resetOSDFilter = function() {
    try {
        var startTimeInput = document.getElementById("osd-start-time");
        var endTimeInput = document.getElementById("osd-end-time");
        var minLatencyInput = document.getElementById("osd-min-latency");
        var maxLatencyInput = document.getElementById("osd-max-latency");
        
        if (startTimeInput) startTimeInput.value = "";
        if (endTimeInput) endTimeInput.value = "";
        if (minLatencyInput) minLatencyInput.value = "";
        if (maxLatencyInput) maxLatencyInput.value = "";
        
        filterOSDTable();
    } catch (e) {
        console.error("Error in resetOSDFilter:", e);
    }
};

// Transaction Table Filter
window.filterTransactionTable = function() {
    try {
        var startTime = document.getElementById("transaction-start-time").value;
        var endTime = document.getElementById("transaction-end-time").value;
        var minDuration = document.getElementById("transaction-min-duration").value;
        var maxDuration = document.getElementById("transaction-max-duration").value;
        var table = document.getElementById("transaction-table");
        if (!table) return;
        var tr = table.getElementsByTagName("tr");
        
        for (var i = 1; i < tr.length; i++) {
            var tds = tr[i].getElementsByTagName("td");
            if (tds.length < 7) continue; // Skip rows with insufficient cells
            
            var tdStartTime = tds[1].textContent;
            var tdEndTime = tds[5].textContent;
            var tdDuration = parseFloat(tds[6].textContent);
            
            var match = true;
            
            if (startTime) {
                var startDate = new Date(startTime.replace('T', ' '));
                var rowStartDate = new Date(tdStartTime);
                if (rowStartDate < startDate) match = false;
            }
            
            if (endTime) {
                var endDate = new Date(endTime.replace('T', ' '));
                var rowEndDate = new Date(tdEndTime);
                if (rowEndDate > endDate) match = false;
            }
            
            if (minDuration) {
                if (tdDuration < parseFloat(minDuration)) match = false;
            }
            
            if (maxDuration) {
                if (tdDuration > parseFloat(maxDuration)) match = false;
            }
            
            tr[i].style.display = match ? "" : "none";
        }

        window.initPagination("transaction-table");
    } catch (e) {
        console.error("Error in filterTransactionTable:", e);
    }
};

window.resetTransactionFilter = function() {
    try {
        var startTimeInput = document.getElementById("transaction-start-time");
        var endTimeInput = document.getElementById("transaction-end-time");
        var minDurationInput = document.getElementById("transaction-min-duration");
        var maxDurationInput = document.getElementById("transaction-max-duration");
        
        if (startTimeInput) startTimeInput.value = "";
        if (endTimeInput) endTimeInput.value = "";
        if (minDurationInput) minDurationInput.value = "";
        if (maxDurationInput) maxDurationInput.value = "";
        
        filterTransactionTable();
    } catch (e) {
        console.error("Error in resetTransactionFilter:", e);
    }
};

// Metadata Sync Table Filter
window.filterMetadataTable = function() {
    try {
        var startTime = document.getElementById("metadata-start-time").value;
        var endTime = document.getElementById("metadata-end-time").value;
        var minDuration = document.getElementById("metadata-min-duration").value;
        var maxDuration = document.getElementById("metadata-max-duration").value;
        var table = document.getElementById("metadata-table");
        if (!table) return;
        var tr = table.getElementsByTagName("tr");
        
        for (var i = 1; i < tr.length; i++) {
            var tdTime = tr[i].getElementsByTagName("td")[0].textContent;
            var tdDuration = parseFloat(tr[i].getElementsByTagName("td")[3].textContent);
            
            var match = true;
            
            if (startTime) {
                var startDate = new Date(startTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate < startDate) match = false;
            }
            
            if (endTime) {
                var endDate = new Date(endTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate > endDate) match = false;
            }
            
            if (minDuration) {
                if (tdDuration < parseFloat(minDuration)) match = false;
            }
            
            if (maxDuration) {
                if (tdDuration > parseFloat(maxDuration)) match = false;
            }
            
            tr[i].style.display = match ? "" : "none";
        }

        window.initPagination("metadata-table");
    } catch (e) {
        console.error("Error in filterMetadataTable:", e);
    }
};

window.resetMetadataFilter = function() {
    try {
        var startTimeInput = document.getElementById("metadata-start-time");
        var endTimeInput = document.getElementById("metadata-end-time");
        var minDurationInput = document.getElementById("metadata-min-duration");
        var maxDurationInput = document.getElementById("metadata-max-duration");
        
        if (startTimeInput) startTimeInput.value = "";
        if (endTimeInput) endTimeInput.value = "";
        if (minDurationInput) minDurationInput.value = "";
        if (maxDurationInput) maxDurationInput.value = "";
        
        filterMetadataTable();
    } catch (e) {
        console.error("Error in resetMetadataFilter:", e);
    }
};

// Client Table Filter
window.filterClientTable = function() {
    try {
        var startTime = document.getElementById("client-start-time").value;
        var endTime = document.getElementById("client-end-time").value;
        var minLatency = document.getElementById("client-min-latency").value;
        var maxLatency = document.getElementById("client-max-latency").value;
        var table = document.getElementById("client-table");
        if (!table) return;
        var tr = table.getElementsByTagName("tr");
        
        for (var i = 1; i < tr.length; i++) {
            var tdTime = tr[i].getElementsByTagName("td")[0].textContent;
            var tdLatency = parseFloat(tr[i].getElementsByTagName("td")[6].textContent);
            
            var match = true;
            
            if (startTime) {
                var startDate = new Date(startTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate < startDate) match = false;
            }
            
            if (endTime) {
                var endDate = new Date(endTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate > endDate) match = false;
            }
            
            if (minLatency) {
                if (tdLatency < parseFloat(minLatency)) match = false;
            }
            
            if (maxLatency) {
                if (tdLatency > parseFloat(maxLatency)) match = false;
            }
            
            tr[i].style.display = match ? "" : "none";
        }

        window.initPagination("client-table");
    } catch (e) {
        console.error("Error in filterClientTable:", e);
    }
};

window.resetClientFilter = function() {
    try {
        var startTimeInput = document.getElementById("client-start-time");
        var endTimeInput = document.getElementById("client-end-time");
        var minLatencyInput = document.getElementById("client-min-latency");
        var maxLatencyInput = document.getElementById("client-max-latency");
        
        if (startTimeInput) startTimeInput.value = "";
        if (endTimeInput) endTimeInput.value = "";
        if (minLatencyInput) minLatencyInput.value = "";
        if (maxLatencyInput) maxLatencyInput.value = "";
        
        filterClientTable();
    } catch (e) {
        console.error("Error in resetClientFilter:", e);
    }
};

// Client Detail Table Filter
window.filterClientDetailTable = function() {
    try {
        var startTime = document.getElementById("client-detail-start-time").value;
        var endTime = document.getElementById("client-detail-end-time").value;
        var minLatency = document.getElementById("client-detail-min-latency").value;
        var maxLatency = document.getElementById("client-detail-max-latency").value;
        var table = document.getElementById("client-detail-table");
        if (!table) return;
        var tr = table.getElementsByTagName("tr");
        
        for (var i = 1; i < tr.length; i++) {
            var tdTime = tr[i].getElementsByTagName("td")[0].textContent;
            var tdLatency = parseFloat(tr[i].getElementsByTagName("td")[6].textContent);
            
            var match = true;
            
            if (startTime) {
                var startDate = new Date(startTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate < startDate) match = false;
            }
            
            if (endTime) {
                var endDate = new Date(endTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate > endDate) match = false;
            }
            
            if (minLatency) {
                if (tdLatency < parseFloat(minLatency)) match = false;
            }
            
            if (maxLatency) {
                if (tdLatency > parseFloat(maxLatency)) match = false;
            }
            
            tr[i].style.display = match ? "" : "none";
        }

        window.initPagination("client-detail-table");
    } catch (e) {
        console.error("Error in filterClientDetailTable:", e);
    }
};

window.resetClientDetailFilter = function() {
    try {
        var startTimeInput = document.getElementById("client-detail-start-time");
        var endTimeInput = document.getElementById("client-detail-end-time");
        var minLatencyInput = document.getElementById("client-detail-min-latency");
        var maxLatencyInput = document.getElementById("client-detail-max-latency");

        if (startTimeInput) startTimeInput.value = "";
        if (endTimeInput) endTimeInput.value = "";
        if (minLatencyInput) minLatencyInput.value = "";
        if (maxLatencyInput) maxLatencyInput.value = "";

        filterClientDetailTable();
    } catch (e) {
        console.error("Error in resetClientDetailFilter:", e);
    }
};

// Dequeue Table Filter
window.filterDequeueTable = function() {
    try {
        var startTime = document.getElementById("dequeue-start-time").value;
        var endTime = document.getElementById("dequeue-end-time").value;
        var minLatency = document.getElementById("dequeue-min-latency").value;
        var maxLatency = document.getElementById("dequeue-max-latency").value;
        var opType = document.getElementById("dequeue-op-type").value;
        var table = document.getElementById("dequeue-table");
        if (!table) return;
        var tr = table.getElementsByTagName("tr");
        
        console.log("Filtering dequeue table with opType:", opType);
        
        for (var i = 1; i < tr.length; i++) {
            var tds = tr[i].getElementsByTagName("td");
            if (tds.length < 6) continue;
            
            var tdTime = tds[0].textContent.trim();
            var tdOpType = tds[1].textContent.trim();
            var tdLatency = parseFloat(tds[5].textContent.trim()); // 现在是毫秒
            
            console.log("Row " + i + " opType:", tdOpType);
            
            var match = true;
            
            if (startTime) {
                var startDate = new Date(startTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate < startDate) match = false;
            }
            
            if (endTime) {
                var endDate = new Date(endTime.replace('T', ' '));
                var rowDate = new Date(tdTime);
                if (rowDate > endDate) match = false;
            }
            
            if (minLatency) {
                if (tdLatency < parseFloat(minLatency)) match = false;
            }
            
            if (maxLatency) {
                if (tdLatency > parseFloat(maxLatency)) match = false;
            }
            
            if (opType) {
                if (tdOpType !== opType) match = false;
            }
            
            tr[i].style.display = match ? "" : "none";
        }

        window.initPagination("dequeue-table");
    } catch (e) {
        console.error("Error in filterDequeueTable:", e);
    }
};

window.resetDequeueFilter = function() {
    try {
        var startTimeInput = document.getElementById("dequeue-start-time");
        var endTimeInput = document.getElementById("dequeue-end-time");
        var minLatencyInput = document.getElementById("dequeue-min-latency");
        var maxLatencyInput = document.getElementById("dequeue-max-latency");
        var opTypeInput = document.getElementById("dequeue-op-type");
        
        if (startTimeInput) startTimeInput.value = "";
        if (endTimeInput) endTimeInput.value = "";
        if (minLatencyInput) minLatencyInput.value = "";
        if (maxLatencyInput) maxLatencyInput.value = "";
        if (opTypeInput) opTypeInput.value = "";
        
        window.filterDequeueTable();
    } catch (e) {
        console.error("Error in resetDequeueFilter:", e);
    }
};

// OSD Op Table Filter
window.filterOSDOpTable = function() {
    try {
        var startTime = document.getElementById("osdop-start-time").value;
        var endTime = document.getElementById("osdop-end-time").value;
        var minDuration = document.getElementById("osdop-min-duration").value;
        var maxDuration = document.getElementById("osdop-max-duration").value;
        var table = document.getElementById("osdop-table");
        if (!table) return;
        var tr = table.getElementsByTagName("tr");
        
        for (var i = 1; i < tr.length; i++) {
            var tdStartTime = tr[i].getElementsByTagName("td")[0].textContent;
            var tdEndTime = tr[i].getElementsByTagName("td")[1].textContent;
            var tdDuration = parseFloat(tr[i].getElementsByTagName("td")[2].textContent);
            
            var match = true;
            
            if (startTime) {
                var startDate = new Date(startTime.replace('T', ' '));
                var rowStartDate = new Date(tdStartTime);
                if (rowStartDate < startDate) match = false;
            }
            
            if (endTime) {
                var endDate = new Date(endTime.replace('T', ' '));
                var rowEndDate = new Date(tdEndTime);
                if (rowEndDate > endDate) match = false;
            }
            
            if (minDuration) {
                if (tdDuration < parseFloat(minDuration)) match = false;
            }
            
            if (maxDuration) {
                if (tdDuration > parseFloat(maxDuration)) match = false;
            }
            
            tr[i].style.display = match ? "" : "none";
        }

        window.initPagination("osdop-table");
    } catch (e) {
        console.error("Error in filterOSDOpTable:", e);
    }
};

window.resetOSDOpFilter = function() {
    try {
        var startTimeInput = document.getElementById("osdop-start-time");
        var endTimeInput = document.getElementById("osdop-end-time");
        var minDurationInput = document.getElementById("osdop-min-duration");
        var maxDurationInput = document.getElementById("osdop-max-duration");
        
        if (startTimeInput) startTimeInput.value = "";
        if (endTimeInput) endTimeInput.value = "";
        if (minDurationInput) minDurationInput.value = "";
        if (maxDurationInput) maxDurationInput.value = "";
        
        window.filterOSDOpTable();
    } catch (e) {
        console.error("Error in resetOSDOpFilter:", e);
    }
};

// Pagination functionality
window.DEFAULT_PAGE_SIZE = 100;
window.paginationState = {};

window.initPagination = function(tableId, pageSize) {
    try {
        var table = document.getElementById(tableId);
        if (!table) return;

        var rows = table.getElementsByTagName("tr");
        var totalRows = 0;

        for (var i = 1; i < rows.length; i++) {
            if (rows[i].style.display !== "none") {
                totalRows++;
            }
        }

        var size = pageSize || window.DEFAULT_PAGE_SIZE;
        if (totalRows === 0) {
            totalRows = rows.length - 1;
        }

        window.paginationState[tableId] = {
            currentPage: 1,
            pageSize: size,
            totalRows: totalRows,
            totalPages: Math.ceil(totalRows / size)
        };

        window.showPage(tableId, 1);
    } catch (e) {
        console.error("Error in initPagination:", e);
    }
};

window.showPage = function(tableId, pageNum) {
    try {
        var state = window.paginationState[tableId];
        if (!state) return;

        var table = document.getElementById(tableId);
        if (!table) return;

        var rows = table.getElementsByTagName("tr");
        var visibleRows = [];

        for (var i = 1; i < rows.length; i++) {
            if (rows[i].style.display !== "none") {
                visibleRows.push(rows[i]);
            }
        }

        var startIndex = (pageNum - 1) * state.pageSize;
        var endIndex = Math.min(startIndex + state.pageSize, visibleRows.length);

        for (var i = 0; i < visibleRows.length; i++) {
            if (i >= startIndex && i < endIndex) {
                visibleRows[i].style.display = "";
            } else {
                visibleRows[i].style.display = "none";
            }
        }

        state.currentPage = pageNum;
        window.updatePaginationControls(tableId);
    } catch (e) {
        console.error("Error in showPage:", e);
    }
};

window.updatePaginationControls = function(tableId) {
    try {
        var state = window.paginationState[tableId];
        if (!state) return;

        var paginationInfo = document.getElementById(tableId + "-pagination-info");
        var paginationButtons = document.getElementById(tableId + "-pagination-buttons");

        if (paginationInfo) {
            var startRow = (state.currentPage - 1) * state.pageSize + 1;
            var endRow = Math.min(state.currentPage * state.pageSize, state.totalRows);
            paginationInfo.textContent = startRow + "-" + endRow + " / " + state.totalRows;
        }

        if (paginationButtons) {
            var html = "";
            var maxButtons = 5;
            var startPage = Math.max(1, state.currentPage - Math.floor(maxButtons / 2));
            var endPage = Math.min(state.totalPages, startPage + maxButtons - 1);

            if (startPage > 1) {
                html += '<button onclick="showPage(\'' + tableId + '\', 1)">&laquo;</button> ';
                html += '<button onclick="showPage(\'' + tableId + '\', ' + (state.currentPage - 1) + ')">&lsaquo;</button> ';
            }

            for (var i = startPage; i <= endPage; i++) {
                if (i === state.currentPage) {
                    html += '<button class="active">' + i + '</button> ';
                } else {
                    html += '<button onclick="showPage(\'' + tableId + '\', ' + i + ')">' + i + '</button> ';
                }
            }

            if (endPage < state.totalPages) {
                html += '<button onclick="showPage(\'' + tableId + '\', ' + (state.currentPage + 1) + ')">&rsaquo;</button> ';
                html += '<button onclick="showPage(\'' + tableId + '\', ' + state.totalPages + ')">&raquo;</button> ';
            }

            paginationButtons.innerHTML = html;
        }
    } catch (e) {
        console.error("Error in updatePaginationControls:", e);
    }
};

window.changePageSize = function(tableId, newSize) {
    try {
        var state = window.paginationState[tableId];
        if (!state) return;

        state.pageSize = parseInt(newSize);
        state.totalPages = Math.ceil(state.totalRows / state.pageSize);
        state.currentPage = 1;

        window.showPage(tableId, 1);
    } catch (e) {
        console.error("Error in changePageSize:", e);
    }
};

window.applyFilterAndPaginate = function(tableId) {
    try {
        var filterFuncName = tableId.replace("-table", "") + "Table";
        if (typeof window[filterFuncName] === 'function') {
            window[filterFuncName]();
        }

        setTimeout(function() {
            window.initPagination(tableId);
        }, 100);
    } catch (e) {
        console.error("Error in applyFilterAndPaginate:", e);
    }
};

document.addEventListener("DOMContentLoaded", function() {
    setTimeout(function() {
        window.initPagination("aio-table");
        window.initPagination("repop-table");
        window.initPagination("osdop-table");
        window.initPagination("osd-table");
        window.initPagination("transaction-table");
        window.initPagination("metadata-table");
        window.initPagination("client-table");
        window.initPagination("client-detail-table");
        window.initPagination("dequeue-table");
    }, 500);
});
