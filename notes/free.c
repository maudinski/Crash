//this logic is wrong, it's not checking if the memBlock is in use
// it'll be something like this i guess
// also, these things wont be in a struct, ill have to access the values depending on
// how i structure them, example: if the first thing stored at the address is the size,
// then the address is the size, of the second thing is the next pointer, then then
// the adress + 8 is the next pointer (8 is the size of an int)

void free(*MemBlock freeThis) {
      MemBlock *prev = freeThis->prev;
      // check previous for either consolidate or link in
      if (freeThis->prev->size == prev->addr - freeThis->addr) { // if theyre touching
            prev->size = prev->size + freeThis->size; // consolidate
            freeThis = freeThis->prev; // set for the algorithm to work
      } else {
            freeThis->next = prev->next;
            prev->next = freeThis;
      }
      MemBlock *next = freeThis->next
      if (freeThis->size == freeThis->addr - next->addr) {
            freeThis->size = freeThis->size + next->size;
            freeThis->next = next->next;
      }
      freeThis->free = true
}