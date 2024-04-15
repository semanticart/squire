(function () {
  var SQUIRE = {
    clear: function () {
      var activeEls = document.getElementsByClassName("active");
      while (activeEls.length) {
        activeEls[0].classList.remove("active");
      }
      setTimeout(function () {
        window.scrollTo(0, 0);
      });
    },

    showSection: function (el) {
      el.classList.add("active");

      while (el.nextElementSibling && el.nextElementSibling.tagName !== "H1") {
        el.nextElementSibling.classList.add("active");
        el = el.nextElementSibling;
      }
    },

    visit: function () {
      var id = window.location.hash.split("#")[1] || "intro";
      var el = document.getElementById(id);
      this.clear();
      this.showSection(el);
    },
  };

  SQUIRE.visit();

  window.addEventListener("hashchange", SQUIRE.visit.bind(SQUIRE));
})();
