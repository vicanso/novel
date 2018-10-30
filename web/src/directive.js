import Vue from "vue";

const intersecting = "_intersecting";
Vue.directive("intersecting-show", {
  bind: function(el, binding) {
    if (!binding.value) {
      return;
    }
    const io = new IntersectionObserver(entries => {
      const target = entries[0];
      // 在元素可见时加载图标，并做diconnect
      if (target.isIntersecting) {
        io.disconnect();
        binding.value();
        delete el.dataset[intersecting];
      }
    });
    io.observe(el);
    el.dataset[intersecting] = io;
  },
  unbind: function(el) {
    const io = el.dataset[intersecting];
    if (io && io.disconnect) {
      io.disconnect();
    }
  }
});
