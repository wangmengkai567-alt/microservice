<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue';

interface Props {
  isOpen: boolean;
  title?: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  closeOnEscape?: boolean;
  closeOnBackdrop?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  closeOnEscape: true,
  closeOnBackdrop: true,
});

const emit = defineEmits<{
  close: [];
}>();

const modalRef = ref<HTMLElement | null>(null);
const previousActiveElement = ref<HTMLElement | null>(null);

// Size classes
const sizeClasses = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
};

// Handle ESC key
const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && props.closeOnEscape && props.isOpen) {
    emit('close');
  }
};

// Handle backdrop click
const handleBackdropClick = (event: MouseEvent) => {
  if (props.closeOnBackdrop && event.target === event.currentTarget) {
    emit('close');
  }
};

// Focus management
const manageFocus = () => {
  if (props.isOpen) {
    // Save current focus
    previousActiveElement.value = document.activeElement as HTMLElement;

    // Focus first focusable element in modal
    setTimeout(() => {
      const focusableElements = modalRef.value?.querySelectorAll(
        'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
      );
      if (focusableElements && focusableElements.length > 0) {
        (focusableElements[0] as HTMLElement).focus();
      }
    }, 100);
  } else {
    // Restore previous focus
    if (previousActiveElement.value) {
      previousActiveElement.value.focus();
    }
  }
};

// Trap focus within modal
const trapFocus = (event: KeyboardEvent) => {
  if (!props.isOpen || !modalRef.value) return;

  const focusableElements = modalRef.value.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  );

  if (focusableElements.length === 0) return;

  const firstElement = focusableElements[0] as HTMLElement;
  const lastElement = focusableElements[focusableElements.length - 1] as HTMLElement;

  if (event.key === 'Tab') {
    if (event.shiftKey) {
      // Shift + Tab
      if (document.activeElement === firstElement) {
        event.preventDefault();
        lastElement.focus();
      }
    } else {
      // Tab
      if (document.activeElement === lastElement) {
        event.preventDefault();
        firstElement.focus();
      }
    }
  }
};

// Watch for open state changes
watch(
  () => props.isOpen,
  (isOpen) => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
      manageFocus();
    } else {
      document.body.style.overflow = '';
      manageFocus();
    }
  }
);

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
  document.addEventListener('keydown', trapFocus);
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown);
  document.removeEventListener('keydown', trapFocus);
  document.body.style.overflow = '';
});
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 overflow-y-auto"
        aria-labelledby="modal-title"
        role="dialog"
        aria-modal="true"
      >
        <!-- Backdrop -->
        <div
          class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
          @click="handleBackdropClick"
        ></div>

        <!-- Modal container -->
        <div class="flex min-h-full items-center justify-center p-4">
          <div
            ref="modalRef"
            :class="[
              'relative transform overflow-hidden rounded-lg bg-white shadow-xl transition-all w-full',
              sizeClasses[size],
            ]"
          >
            <!-- Header -->
            <div v-if="title || $slots.header" class="bg-white px-4 pt-5 pb-4 sm:p-6">
              <div class="flex items-start justify-between">
                <h3
                  v-if="title"
                  id="modal-title"
                  class="text-lg font-medium leading-6 text-gray-900"
                >
                  {{ title }}
                </h3>
                <slot name="header" />
                <button
                  type="button"
                  class="ml-auto rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                  @click="emit('close')"
                  aria-label="关闭对话框"
                >
                  <span class="sr-only">关闭</span>
                  <svg
                    class="h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    aria-hidden="true"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              </div>
            </div>

            <!-- Body -->
            <div class="bg-white px-4 pb-4 sm:p-6">
              <slot />
            </div>

            <!-- Footer -->
            <div
              v-if="$slots.footer"
              class="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6"
            >
              <slot name="footer" />
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .relative,
.modal-leave-active .relative {
  transition: transform 0.3s ease;
}

.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: scale(0.95);
}
</style>
